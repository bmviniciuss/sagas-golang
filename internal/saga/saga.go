package saga

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ActionType string

func (at ActionType) String() string {
	return string(at)
}

func (at ActionType) IsSuccess() bool {
	return at == SuccessActionType
}

func (at ActionType) IsFailure() bool {
	return at == FailureActionType
}

func (at ActionType) IsCompensated() bool {
	return at == CompensatedActionType
}

const (
	RequestActionType      ActionType = "request"
	CompensationActionType ActionType = "compensation"
	SuccessActionType      ActionType = "success"
	FailureActionType      ActionType = "failure"
	CompensatedActionType  ActionType = "compensated"
)

type Workflow struct {
	ID           uuid.UUID
	Name         string
	Steps        *StepDoubleLinkedList
	ReplyChannel string
}

func (w Workflow) EventTypes() map[string]uuid.UUID {
	ts := map[string]uuid.UUID{}
	steps := w.Steps.ToList()
	for _, step := range steps {
		ts[fmt.Sprintf("%s.%s.success", w.Name, step.Name)] = w.ID
		ts[fmt.Sprintf("%s.%s.failure", w.Name, step.Name)] = w.ID
	}
	return ts
}

type StepData struct {
	ID             uuid.UUID
	Name           string
	ServiceName    string
	Compensable    bool
	PayloadBuilder func(ctx context.Context, data any) (map[string]interface{}, error)
}

type Step struct {
	*StepData
	next *Step
	prev *Step
}

func (s Step) Next() (*Step, bool) {
	if s.next == nil {
		return nil, false
	}
	return s.next, true
}

func (s Step) Previous() (*Step, bool) {
	if s.prev == nil {
		return nil, false
	}
	return s.prev, true
}

type StepDoubleLinkedList struct {
	head *Step
	tail *Step
	len  int
}

func NewStepList() *StepDoubleLinkedList {
	return &StepDoubleLinkedList{}
}

func (sl *StepDoubleLinkedList) Append(data *StepData) {
	newNode := &Step{
		StepData: data,
		next:     nil,
		prev:     nil,
	}
	if sl.head == nil {
		sl.head = newNode
		sl.tail = newNode
		sl.len++
		return
	}

	newNode.prev = sl.tail
	sl.tail.next = newNode
	sl.tail = newNode
	sl.len++
}

func (sl *StepDoubleLinkedList) Head() (*Step, bool) {
	return sl.head, sl.head != nil
}

func (sl *StepDoubleLinkedList) GetStep(id uuid.UUID) (*Step, bool) {
	current := sl.head
	for current != nil {
		if current.ID == id {
			return current, true
		}
		current = current.next
	}
	return nil, false
}

func (sl *StepDoubleLinkedList) ToList() []Step {
	s := make([]Step, sl.Len())
	current := sl.head
	i := 0
	for current != nil {
		s[i] = *current
		current = current.next
		i++
	}
	return s
}

func (sl *StepDoubleLinkedList) Len() int {
	return sl.len
}
