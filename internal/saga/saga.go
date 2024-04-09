package saga

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Workflow struct {
	ID           uuid.UUID
	Name         string
	Steps        *StepsList
	ReplyChannel string
}

func (w *Workflow) EventTypes() map[string]uuid.UUID {
	ts := map[string]uuid.UUID{}
	steps := w.Steps.ToList()
	for _, step := range steps {
		ts[fmt.Sprintf("%s.%s.success", w.Name, step.Name)] = w.ID
		ts[fmt.Sprintf("%s.%s.failure", w.Name, step.Name)] = w.ID
	}
	return ts
}

var (
	ErrCurrentStepNotFound = fmt.Errorf("current step not found in message workflow")
	ErrUnknownActionType   = fmt.Errorf("unknown action type")
)

// GetNextStep returns the next step in the workflow based on the message received and the current workflow
// If the message is a success message, the next step in the workflow is returned or nil if there are no more steps
// If the message is a failure message, the first compensation step is returned or nil if there are no more steps
// If the message is a compensated message, the next compensable step in the workflow is returned or nil if there are no more steps
func (w *Workflow) GetNextStep(ctx context.Context, message Message) (*Step, error) {
	currentStep, ok := w.Steps.GetStep(message.Saga.Step.ID)
	if !ok {
		return nil, ErrCurrentStepNotFound
	}

	if message.ActionType.IsSuccess() {
		nextStep, ok := currentStep.Next()
		if !ok {
			return nil, nil
		}
		return nextStep, nil
	}

	if message.ActionType.IsFailure() {
		firstCompensableStep, ok := currentStep.FirstCompensableStep()
		if !ok {
			return nil, nil
		}
		return firstCompensableStep, nil
	}

	if message.ActionType.IsCompensated() {
		nextCompensableStep, ok := currentStep.FirstCompensableStep()
		if !ok {
			return nil, nil
		}
		return nextCompensableStep, nil
	}
	return nil, ErrUnknownActionType
}
