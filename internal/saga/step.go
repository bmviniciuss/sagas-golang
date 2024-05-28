package saga

import (
	"context"
	"fmt"
)

// PayloadBuilder is an interface for building the payload of a step in the workflow.
type PayloadBuilder interface {
	// Build builds the payload of a step in the workflow based on the given data and action.
	Build(ctx context.Context, exec *Execution, action ActionType) (map[string]interface{}, error)
}

type (
	// StepData represents the data of a step in the workflow.
	StepData struct {
		Name           string
		ServiceName    string
		Compensable    bool
		PayloadBuilder PayloadBuilder
		EventTypes     EventTypes
		Topics         Topics
	}

	EventTypes struct {
		Request            string
		Success            string
		Failure            string
		CompesationRequest string
		Compensation       string
	}

	Topics struct {
		Request  string
		Response string
	}

	// Step represents a step in the workflow.
	Step struct {
		*StepData
		next *Step
		prev *Step
	}
)

func (sd *StepData) IsSuccess(eventType string) bool {
	return eventType == sd.EventTypes.Success
}

func (sd *StepData) IsFailure(eventType string) bool {
	return eventType == sd.EventTypes.Failure
}

func (sd *StepData) IsCompensation(eventType string) bool {
	return eventType == sd.EventTypes.Compensation
}

// Next returns the next step in the workflow.
//
// Returns the next step if it exists
// returns nil, false if there is no next step
func (s *Step) Next() (*Step, bool) {
	return s.next, s.next != nil
}

// Previous returns the previous step in the workflow.
//
// Returns the previous step if it exists
// returns nil, false if there is no previous step
func (s *Step) Previous() (*Step, bool) {
	return s.prev, s.prev != nil
}

// FirstCompensableStep returns the first compensable step in the workflow before the current step or the current step itself.
//
// returns the current step if it is compensable
//
// returns the first compensable step before the current step
//
// returns nil if no compensable step is found
func (s *Step) FirstCompensableStep() (*Step, bool) {
	prev, ok := s.Previous()
	if !ok {
		return nil, false
	}

	current := prev
	for current != nil {
		if current.Compensable {
			return current, true
		}
		current, _ = current.Previous()
	}
	return nil, false
}

func (s *Step) DestinationTopic(actionType ActionType) string {
	return fmt.Sprintf("service.%s.%s", s.ServiceName, actionType.String())
}

type StepsList struct {
	head *Step
	tail *Step
	len  int
}

// TODO: prevent duplicated step names
func NewStepList(stepsData ...*StepData) *StepsList {
	list := &StepsList{}
	for _, data := range stepsData {
		list.Append(data)
	}
	return list
}

// Append adds a new step to the workflow.
// It returns the newly added step.
func (sl *StepsList) Append(data *StepData) *Step {
	newNode := &Step{
		StepData: data,
		next:     nil,
		prev:     nil,
	}
	if sl.head == nil {
		sl.head = newNode
		sl.tail = newNode
		sl.len++
		return newNode
	}

	newNode.prev = sl.tail
	sl.tail.next = newNode
	sl.tail = newNode
	sl.len++
	return newNode
}

// Head returns the first step in the workflow.
func (sl *StepsList) Head() (*Step, bool) {
	return sl.head, sl.head != nil
}

// GetStep returns the step with the given id.
// TODO: remove this method
func (sl *StepsList) GetStep(name string) (*Step, bool) {
	current := sl.head
	for current != nil {
		if current.Name == name {
			return current, true
		}
		current, _ = current.Next()
	}
	return nil, false
}

// GetStep returns the step with the given id.
func (sl *StepsList) GetStepFromServiceEvent(serviceName string, eventType string) (*Step, bool) {
	current := sl.head
	for current != nil {
		if current.ServiceName == serviceName {
			if current.EventTypes.Request == eventType {
				return current, true
			}
			if current.EventTypes.Success == eventType {
				return current, true
			}
			if current.EventTypes.Failure == eventType {
				return current, true
			}
			if current.EventTypes.Compensation == eventType {
				return current, true
			}
		}
		current, _ = current.Next()
	}
	return nil, false
}

// ToList returns the steps in the workflow as a slice.
func (sl *StepsList) ToList() []Step {
	s := make([]Step, sl.Len())
	current := sl.head
	i := 0
	for current != nil {
		s[i] = *current
		current, _ = current.Next()
		i++
	}
	return s
}

// Len returns the number of steps in the workflow.
func (sl *StepsList) Len() int {
	return sl.len
}
