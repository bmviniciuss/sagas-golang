package saga

import (
	"context"
	"fmt"
	"reflect"

	"github.com/bmviniciuss/sagas-golang/pkg/events"
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
func (w *Workflow) GetNextStep(ctx context.Context, message events.Event) (*Step, error) { // TODO: use ptr to workflow
	currentStep, ok := w.Steps.GetStepFromServiceEvent(message.Origin, message.Type)
	if !ok {
		return nil, ErrCurrentStepNotFound
	}

	if currentStep.IsSuccess(message.Type) {
		nextStep, ok := currentStep.Next()
		if !ok {
			return nil, nil
		}
		return nextStep, nil
	}

	if currentStep.IsFailure(message.Type) {
		firstCompensableStep, ok := currentStep.FirstCompensableStep()
		if !ok {
			return nil, nil
		}
		return firstCompensableStep, nil
	}

	if currentStep.IsCompensation(message.Type) {
		nextCompensableStep, ok := currentStep.FirstCompensableStep()
		if !ok {
			return nil, nil
		}
		return nextCompensableStep, nil
	}
	return nil, ErrUnknownActionType
}
