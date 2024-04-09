package saga

import (
	"context"
	"fmt"
)

type Coordinator struct {
}

func NewCoordinator() *Coordinator {
	return &Coordinator{}
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
	fmt.Println("Processing message")
	currentStep, ok := workflow.Steps.GetStep(message.Saga.Step.ID)
	if !ok {
		return nil, ErrCurrentStepNotFound
	}

	if message.ActionType.IsSuccess() {
		fmt.Println("Action Type is success. Proceeding to saga's next step")
		nextStep, ok := currentStep.Next()
		if !ok {
			fmt.Println("Saga does not have next step and has finished")
			return nil, nil
		}
		fmt.Println("Saga has next step")
		return nextStep, nil
	}

	if message.ActionType.IsFailure() {
		fmt.Println("Action Type is failure. Start compensation flow from the first compensable step")
		firstCompensableStep, ok := currentStep.FirstCompensableStep()
		if !ok {
			fmt.Println("No compensable steps found")
			return nil, nil
		}
		fmt.Printf("First compensable step found with id [%s] and name [%s]\n", firstCompensableStep.ID.String(), firstCompensableStep.Name)
		return firstCompensableStep, nil
	}

	if message.ActionType.IsCompensated() {
		fmt.Println("Action Type is compensated. Proceeding to saga's next compensable step")
		nextCompensableStep, ok := currentStep.FirstCompensableStep()
		if !ok {
			fmt.Println("No compensable steps found")
			return nil, nil
		}
		fmt.Printf("First compensable step found with id [%s] and name [%s]\n", nextCompensableStep.ID.String(), nextCompensableStep.Name)
		return nextCompensableStep, nil
	}

	return nil, ErrUnknownActionType
}
