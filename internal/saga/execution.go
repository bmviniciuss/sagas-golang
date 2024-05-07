package saga

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/bmviniciuss/sagas-golang/pkg/structs"
	"github.com/google/uuid"
)

type Execution struct {
	ID       uuid.UUID
	Workflow *Workflow
	State    map[string][]byte
}

func (e *Execution) IsEmpty() bool {
	return reflect.DeepEqual(e, &Execution{})
}

func NewExecution(workflow *Workflow) *Execution {
	return &Execution{
		ID:       uuid.New(),
		Workflow: workflow,
		State:    make(map[string][]byte),
	}
}

func (e *Execution) SetState(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	e.State[key] = data
	return nil
}

func (e *Execution) Read(key string, dest interface{}) error {
	data, ok := e.State[key]
	if !ok {
		return errors.New("unable to get key value")
	}
	err := structs.FromBytes(data, &dest)
	if err != nil {
		return err
	}
	return nil
}
