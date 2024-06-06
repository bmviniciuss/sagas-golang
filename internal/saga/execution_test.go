package saga

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewExecution(t *testing.T) {
	t.Run("should create a new execution", func(t *testing.T) {
		workflow := &Workflow{}
		execution := NewExecution(workflow)
		assert.NotNil(t, execution)
		assert.NotNil(t, execution.ID)
		assert.Equal(t, workflow, execution.Workflow)
		assert.NotNil(t, execution.State)
	})
}

func TestExecution_SetState(t *testing.T) {
	t.Run("should set state", func(t *testing.T) {
		execution := NewExecution(&Workflow{})
		execution.SetState("key", "value")
		assert.Equal(t, "value", execution.State["key"])
	})
}

func TestExecution_Read(t *testing.T) {
	t.Run("should read state", func(t *testing.T) {
		execution := NewExecution(&Workflow{})
		execution.SetState("key", "value")
		var dest string
		err := execution.Read("key", &dest)
		assert.Nil(t, err)
		assert.Equal(t, "value", dest)
	})

	t.Run("should return error when key does not exist", func(t *testing.T) {
		execution := NewExecution(&Workflow{})
		var dest string
		err := execution.Read("key", &dest)
		assert.Error(t, err)
	})

	t.Run("should return error when value is not a valid type", func(t *testing.T) {
		execution := NewExecution(&Workflow{})
		execution.SetState("key", "value")
		var dest int
		err := execution.Read("key", &dest)
		assert.NotNil(t, err)
	})

	t.Run("should return value if exsits", func(t *testing.T) {
		execution := NewExecution(&Workflow{})
		execution.SetState("key", "value")
		var dest string
		err := execution.Read("key", &dest)
		assert.Nil(t, err)
		assert.Equal(t, "value", dest)
	})
}
