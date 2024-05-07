package service

import (
	"context"
	"encoding/json"

	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Publisher interface {
	Publish(ctx context.Context, destination string, data []byte) error
}

type Port interface {
	ProcessMessage(ctx context.Context, message *saga.Message, execution *saga.Execution) error
	Start(ctx context.Context, workflow *saga.Workflow, data map[string]interface{}) (*uuid.UUID, error)
}

type Execution struct {
	logger              *zap.SugaredLogger
	executionRepository saga.ExecutionRepository
	publisher           Publisher
}

func NewExecution(
	logger *zap.SugaredLogger,
	executionRepository saga.ExecutionRepository,
	publisher Publisher,
) *Execution {
	return &Execution{
		logger:              logger,
		executionRepository: executionRepository,
		publisher:           publisher,
	}
}

func (w *Execution) Start(ctx context.Context, workflow *saga.Workflow, data map[string]interface{}) (*uuid.UUID, error) {
	l := w.logger
	l.Info("Starting workflow")
	execution := saga.NewExecution(workflow)
	l.Infof("Starting saga with ID: %s", execution.ID.String())
	execution.SetState("input", data)
	err := w.executionRepository.Save(ctx, execution)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error while saving execution")
		return nil, err
	}

	firstStep, ok := execution.Workflow.Steps.Head()
	if !ok {
		l.Info("There are no steps to process. Successfully finished workflow.")
		return nil, nil
	}
	actionType := saga.RequestActionType
	payload, err := firstStep.PayloadBuilder.Build(ctx, data, actionType)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error while building payload")
		return nil, err
	}

	firstMsg := saga.NewMessage(execution.ID, payload, nil, workflow, firstStep, actionType)
	jsonMsg, err := firstMsg.ToJSON()
	if err != nil {
		l.With(zap.Error(err)).Error("Got error while marshalling message")
		return nil, err
	}
	err = w.publisher.Publish(ctx, firstStep.DestinationTopic(actionType), jsonMsg)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error publishing message to destination")
		return nil, err
	}
	l.Info("Successfully started workflow")
	return &execution.ID, nil
}

// TODO: add unit tests
func (w *Execution) ProcessMessage(ctx context.Context, message *saga.Message, execution *saga.Execution) (err error) {
	lggr := w.logger
	lggr.Infof("Processing message: %s", message.EventType.Action.String())
	workflow := execution.Workflow
	// execution.SetState(fmt.Sprintf("%s.response", message.Saga.Step.StateKey()), message.EventData)
	err = w.executionRepository.Save(ctx, execution)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error saving execution state")
		return err
	}

	nextStep, err := workflow.GetNextStep(ctx, *message)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error while getting next step")
		return err
	}
	if nextStep == nil {
		lggr.Info("There are no more steps to process. Successfully finished workflow.")
		return nil
	}
	lggr.Infof("Next step: %s", nextStep.Name)
	nextActionType, err := message.EventType.Action.Next()
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error while getting next action type")
		return err
	}
	lggr.Infof("Next step action type: %s", nextStep.Name)
	payload, err := nextStep.PayloadBuilder.Build(ctx, execution.State, nextActionType)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error while building payload")
		return err
	}
	nextMsg := saga.NewMessage(message.GlobalID, payload, message.Metadata, workflow, nextStep, nextActionType)
	// execution.SetState(fmt.Sprintf("%s.request", nextMsg.Saga.Step.StateKey()), nextMsg.EventData)
	err = w.executionRepository.Save(ctx, execution)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error saving execution next step request state")
		return err
	}

	jsonMsg, err := json.Marshal(nextMsg)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error while marshalling message")
		return err
	}

	err = w.publisher.Publish(ctx, nextStep.DestinationTopic(nextActionType), jsonMsg)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error publishing message to destination")
		return err
	}

	lggr.Infof("Successfully processed message and produce")
	return nil
}
