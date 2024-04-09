package saga

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type Coordinator struct {
	logger *zap.SugaredLogger
}

func NewCoordinator(logger *zap.SugaredLogger) *Coordinator {
	return &Coordinator{
		logger: logger,
	}
}

var (
	ErrCurrentStepNotFound = fmt.Errorf("current step not found in message workflow")
	ErrUnknownActionType   = fmt.Errorf("unknown action type")
)

// GetNextStep returns the next step in the workflow based on the message received and the current workflow
// If the message is a success message, the next step in the workflow is returned or nil if there are no more steps
// If the message is a failure message, the first compensation step is returned or nil if there are no more steps
// If the message is a compensated message, the next compensable step in the workflow is returned or nil if there are no more steps
func (c *Coordinator) GetNextStep(ctx context.Context, message Message, workflow Workflow) (*Step, error) {
	l := c.logger
	l.Info("Processing message with event_id [%s] and event_type [%s]", message.EventID, message.EventType)
	currentStep, ok := workflow.Steps.GetStep(message.Saga.Step.ID)
	if !ok {
		l.Error("Current step not found in message's workflow")
		return nil, ErrCurrentStepNotFound
	}

	if message.ActionType.IsSuccess() {
		l.Info("Action Type is success. Proceeding to saga's next step")
		nextStep, ok := currentStep.Next()
		if !ok {
			l.Info("Saga does not have next step and has finished")
			return nil, nil
		}
		l.Info("Saga has next step")
		return nextStep, nil
	}

	if message.ActionType.IsFailure() {
		l.Info("Action Type is failure. Start compensation flow from the first compensable step")
		firstCompensableStep, ok := currentStep.FirstCompensableStep()
		if !ok {
			l.Info("No compensable steps found")
			return nil, nil
		}
		l.Infof("First compensable step found with id [%s] and name [%s]\n", firstCompensableStep.ID.String(), firstCompensableStep.Name)
		return firstCompensableStep, nil
	}

	if message.ActionType.IsCompensated() {
		l.Info("Action Type is compensated. Proceeding to saga's next compensable step")
		nextCompensableStep, ok := currentStep.FirstCompensableStep()
		if !ok {
			l.Info("No compensable steps found")
			return nil, nil
		}
		l.Infof("First compensable step found with id [%s] and name [%s]\n", nextCompensableStep.ID.String(), nextCompensableStep.Name)
		return nextCompensableStep, nil
	}

	return nil, ErrUnknownActionType
}
