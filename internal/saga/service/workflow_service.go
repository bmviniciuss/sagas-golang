package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/bmviniciuss/sagas-golang/pkg/events"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Publisher interface {
	Publish(ctx context.Context, destination string, data []byte) error
}

type Port interface {
	ProcessMessage(ctx context.Context, message *events.Event, execution *saga.Execution) error
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
	lggr := w.logger
	lggr.Info("Starting workflow")
	execution := saga.NewExecution(workflow)
	lggr.Infof("Starting saga with ID: %s", execution.ID.String())
	err := execution.SetState("input", data)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error setting input data to execution")
		return nil, err
	}
	err = w.executionRepository.Insert(ctx, execution)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error while saving execution")
		return nil, err
	}

	firstStep, ok := execution.Workflow.Steps.Head()
	if !ok {
		lggr.Info("There are no steps to process. Successfully finished workflow.")
		return nil, nil
	}
	actionType := saga.REQUEST_ACTION_TYPE
	event, err := firstStep.PayloadBuilder.Build(ctx, execution, actionType)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error while building payload")
		return nil, err
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error while marshalling event data")
		return nil, err
	}

	err = w.publisher.Publish(ctx, firstStep.Topics.Request, eventJSON)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error publishing message to destination")
		return nil, err
	}
	lggr.Info("Successfully started workflow")
	return &execution.ID, nil
}

// TODO: add unit tests
func (w *Execution) ProcessMessage(ctx context.Context, event *events.Event, execution *saga.Execution) (err error) {
	lggr := w.logger
	lggr.Infof("Saga Service started processing message with event: %s", event.Type)
	workflow := execution.Workflow
	currentStep, ok := workflow.Steps.GetStepFromServiceEvent(event.Origin, event.Type)
	if !ok {
		return errors.New("currenct step not found in workflow")
	}

	currenctStepResponseKey := fmt.Sprintf("%s.response.%s", currentStep.Name, event.Type)
	// Saving response data to execution state
	err = execution.SetState(currenctStepResponseKey, event.Data)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error setting message data to execution state")
		return err
	}

	err = w.executionRepository.Save(ctx, execution)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error saving execution state")
		return err
	}

	// Aquring next step
	nextStep, err := workflow.GetNextStep(ctx, currentStep, event.Type)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error while getting next step")
		return err
	}
	if nextStep.Step == nil {
		lggr.Info("There are no more steps to process. Successfully finished workflow.")
		return nil
	}
	lggr.Infof("Next step: %s", nextStep.Step.Name)

	// building next event event
	nextEvent, err := nextStep.Step.PayloadBuilder.Build(ctx, execution, nextStep.ActionType)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error building next step event")
		return err
	}

	eventJSON, err := json.Marshal(nextEvent)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error while marshalling event data")
		return err
	}

	err = w.publisher.Publish(ctx, nextStep.Step.Topics.Request, eventJSON)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error publishing message to destination")
		return err
	}

	lggr.Infof("Successfully processed message and produce")
	return nil
}
