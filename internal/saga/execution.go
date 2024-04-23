package saga

import (
	"errors"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
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

func (e *Execution) SetState(key string, value interface{}) {
	e.State[key] = value
}

func (e *Execution) Read(key string, dest interface{}) error {
	data, ok := e.State[key].(map[string]interface{})
	if !ok {
		return errors.New("unable to get key value")
	}
	err := mapstructure.Decode(data, dest) // TODO: remove this dependency
	if err != nil {
		return err
	}
	return nil
}
