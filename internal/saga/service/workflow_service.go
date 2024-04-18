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
	ProcessMessage(ctx context.Context, message *saga.Message, workflow *saga.Workflow) error
	Start(ctx context.Context, workflow *saga.Workflow, data map[string]interface{}) (*uuid.UUID, error)
}

type Workflow struct {
	logger    *zap.SugaredLogger
	publisher Publisher
}

func NewWorkflow(
	logger *zap.SugaredLogger,
	publisher Publisher,
) *Workflow {
	return &Workflow{
		logger:    logger,
		publisher: publisher,
	}
}

// TODO: add unit tests
func (w *Workflow) ProcessMessage(ctx context.Context, message *saga.Message, workflow *saga.Workflow) error {
	l := w.logger
	l.Infof("Processing message: %s", message.EventType.Action.String())
	nextStep, err := workflow.GetNextStep(ctx, *message)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error while getting next step")
		return err
	}
	if nextStep == nil {
		l.Info("There are no more steps to process. Successfully finished workflow.")
		return nil
	}
	l.Infof("Next step: %s", nextStep.Name)
	nextActionType, err := message.EventType.Action.Next()
	if err != nil {
		l.With(zap.Error(err)).Error("Got error while getting next action type")
		return err
	}
	l.Infof("Next step action type: %s", nextStep.Name)
	payload, err := nextStep.PayloadBuilder.Build(ctx, message.EventData, nextActionType)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error while building payload")
		return err
	}
	nextMsg := saga.NewMessage(message.GlobalID, payload, message.Metadata, workflow, nextStep, nextActionType)
	jsonMsg, err := json.Marshal(nextMsg)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error while marshalling message")
		return err
	}

	err = w.publisher.Publish(ctx, nextStep.DestinationTopic(nextActionType), jsonMsg)
	if err != nil {
		l.With(zap.Error(err)).Error("Got error publishing message to destination")
		return err
	}

	l.Infof("Successfully processed message and produce")
	return nil
}

func (w *Workflow) Start(ctx context.Context, workflow *saga.Workflow, data map[string]interface{}) (*uuid.UUID, error) {
	l := w.logger
	l.Info("Starting workflow")
	firstStep, ok := workflow.Steps.Head()
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
	globalID := uuid.New()
	l.Infof("Starting saga with ID: %s", globalID.String())
	firstMsg := saga.NewMessage(globalID, payload, nil, workflow, firstStep, actionType)
	jsonMsg, err := json.Marshal(firstMsg)
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
	return &globalID, nil
}
