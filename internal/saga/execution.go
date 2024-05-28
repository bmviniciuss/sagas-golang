package saga

import (
	"errors"
	"reflect"

	"github.com/bmviniciuss/sagas-golang/pkg/structs"
	"github.com/google/uuid"
)

type Execution struct {
	ID       uuid.UUID
	Workflow *Workflow
	State    map[string]interface{}
}

func (e *Execution) IsEmpty() bool {
	return reflect.DeepEqual(e, &Execution{})
}

func NewExecution(workflow *Workflow) *Execution {
	return &Execution{
		ID:       uuid.New(),
		Workflow: workflow,
		State:    make(map[string]interface{}),
	}
}

func (e *Execution) SetState(key string, value interface{}) error {
	e.State[key] = value
	return nil
}

func (e *Execution) Read(key string, dest interface{}) error {
	data, ok := e.State[key]
	if !ok {
		return errors.New("unable to get key value")
	}

	dataBytes, err := structs.ToBytes(data)
	if err != nil {
		return err
	}

	err = structs.FromBytes(dataBytes, &dest)
	if err != nil {
		return err
	}
	return nil
}
