package saga

import (
	"context"
	"fmt"
	"reflect"
)

var (
	ErrCurrentStepNotFound = fmt.Errorf("current step not found in message workflow")
	ErrUnknownActionType   = fmt.Errorf("unknown action type")
)

type Workflow struct {
	Name         string
	ReplyChannel string
	Steps        *StepsList
}

// IsEmpty returns true if the workflow is empty
func (w *Workflow) IsEmpty() bool {
	return reflect.DeepEqual(w, &Workflow{})
}

// GetNextStep returns the next step in the workflow based on the message received and the current workflow
// If the message is a success message, the next step in the workflow is returned or nil if there are no more steps
// If the message is a failure message, the first compensation step is returned or nil if there are no more steps
// If the message is a compensated message, the next compensable step in the workflow is returned or nil if there are no more steps
func (w *Workflow) GetNextStep(ctx context.Context, message Message) (*Step, error) { // TODO: use ptr to workflow
	currentStep, ok := w.Steps.GetStep(message.Saga.Step.Name)
	if !ok {
		return nil, ErrCurrentStepNotFound
	}

	if message.EventType.Action.IsSuccess() {
		nextStep, ok := currentStep.Next()
		if !ok {
			return nil, nil
		}
		return nextStep, nil
	}

	if message.EventType.Action.IsFailure() {
		firstCompensableStep, ok := currentStep.FirstCompensableStep()
		if !ok {
			return nil, nil
		}
		return firstCompensableStep, nil
	}

	if message.EventType.Action.IsCompensated() {
		nextCompensableStep, ok := currentStep.FirstCompensableStep()
		if !ok {
			return nil, nil
		}
		return nextCompensableStep, nil
	}
	return nil, ErrUnknownActionType
}
