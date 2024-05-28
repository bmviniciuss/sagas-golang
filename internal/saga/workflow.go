package saga

import (
	"context"
	"fmt"
	"reflect"
)

var (
	ErrUnknownActionType = fmt.Errorf("unknown action type")
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

type NextStep struct {
	Step       *Step
	ActionType ActionType
}

// GetNextStep returns the next step in the workflow based on current step and a received event type
// If the event type is a success message, the next step in the workflow is returned or nil if there are no more steps
// If the event type is a failure message, the first compensation step is returned or nil if there are no more steps
// If the event type is a compensated message, the next compensable step in the workflow is returned or nil if there are no more steps
func (w *Workflow) GetNextStep(ctx context.Context, currentStep *Step, eventType string) (NextStep, error) {
	if currentStep.IsSuccess(eventType) {
		nextStep, ok := currentStep.Next()
		if !ok {
			return NextStep{}, nil
		}
		return NextStep{
			Step:       nextStep,
			ActionType: REQUEST_ACTION_TYPE,
		}, nil
	}

	if currentStep.IsFailure(eventType) {
		firstCompensableStep, ok := currentStep.FirstCompensableStep()
		if !ok {
			return NextStep{}, nil
		}
		return NextStep{
			Step:       firstCompensableStep,
			ActionType: COMPESATION_REQUEST_ACTION_TYPE,
		}, nil
	}

	if currentStep.IsCompensation(eventType) {
		nextCompensableStep, ok := currentStep.FirstCompensableStep()
		if !ok {
			return NextStep{}, nil
		}
		return NextStep{
			Step:       nextCompensableStep,
			ActionType: COMPESATION_REQUEST_ACTION_TYPE,
		}, nil
	}
	return NextStep{}, ErrUnknownActionType
}
