package saga

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bmviniciuss/sagas-golang/pkg/events"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Publisher interface {
	Publish(ctx context.Context, destination string, data []byte) error
}

type ServicePort interface {
	Start(ctx context.Context, workflow *Workflow, data map[string]interface{}) (*uuid.UUID, error)
	ProcessMessage(ctx context.Context, message *events.Event, execution *Execution) error
}

type Service struct {
	logger              *zap.SugaredLogger
	executionRepository ExecutionRepository
	publisher           Publisher
}

var (
	_ ServicePort = (*Service)(nil)
)

func NewService(
	logger *zap.SugaredLogger,
	executionRepository ExecutionRepository,
	publisher Publisher,
) *Service {
	return &Service{
		logger:              logger,
		executionRepository: executionRepository,
		publisher:           publisher,
	}
}

func (service *Service) Start(ctx context.Context, workflow *Workflow, data map[string]interface{}) (*uuid.UUID, error) {
	lggr := service.logger
	lggr.Info("Starting workflow")
	execution := NewExecution(workflow)
	lggr.Infof("Starting saga with ID: %s", execution.ID.String())
	err := execution.SetState("input", data)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error setting input data to execution")
		return nil, err
	}
	err = service.executionRepository.Insert(ctx, execution)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error while saving execution")
		return nil, err
	}

	firstStep, ok := execution.Workflow.Steps.Head()
	if !ok {
		lggr.Info("There are no steps to process. Successfully finished workflow.")
		return nil, nil
	}
	actionType := REQUEST_ACTION_TYPE
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

	err = service.publisher.Publish(ctx, firstStep.Topics.Request, eventJSON)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error publishing message to destination")
		return nil, err
	}
	lggr.Info("Successfully started workflow")
	return &execution.ID, nil
}

// TODO: add unit tests
func (service *Service) ProcessMessage(ctx context.Context, event *events.Event, execution *Execution) (err error) {
	lggr := service.logger
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

	err = service.executionRepository.Save(ctx, execution)
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

	err = service.publisher.Publish(ctx, nextStep.Step.Topics.Request, eventJSON)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error publishing message to destination")
		return err
	}

	lggr.Infof("Successfully processed message and produce")
	return nil
}
