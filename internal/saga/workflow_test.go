package saga

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// import (
// 	"context"
// 	"fmt"
// 	"testing"

// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// )

func TestWorkflow_GetNextStep(t *testing.T) {
	t.Run("should return error when current step is not found in message workflow", func(t *testing.T) {
		steps := NewStepList(&StepData{
			ID:          uuid.New(),
			Name:        "create_order",
			ServiceName: "order",
			Compensable: true,
			PayloadBuilder: func(ctx context.Context, data any) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
		})
		workflow := Workflow{
			ID:           uuid.New(),
			Name:         "create_order_saga",
			ReplyChannel: "kfk.dev.create_order_saga.reply",
			Steps:        steps,
		}
		message := Message{
			GlobalID: uuid.New(),
			EventID:  uuid.New(),
			EventType: EventType{
				SagaName: workflow.Name,
				StepName: "create-order",
				Action:   SuccessActionType,
			},
			Saga: Saga{
				ID:           workflow.ID,
				Name:         workflow.Name,
				ReplyChannel: workflow.ReplyChannel,
				Step: SagaStep{
					ID:     uuid.New(),
					Name:   "create-order",
					Action: SuccessActionType,
				},
			},
			Metadata: map[string]string{
				"client_id": uuid.NewString(),
			},
			EventData: map[string]interface{}{},
		}
		step, err := workflow.GetNextStep(context.Background(), message)
		assert.Nil(t, step)
		assert.Error(t, err)
		assert.Equal(t, ErrCurrentStepNotFound, err)
	})

	t.Run("should return nil when message action type is success and there are no more steps", func(t *testing.T) {
		createOrderStep := &StepData{
			ID:          uuid.New(),
			Name:        "create_order",
			ServiceName: "order",
			Compensable: true,
			PayloadBuilder: func(ctx context.Context, data any) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
		}
		steps := NewStepList(createOrderStep)
		workflow := Workflow{
			ID:           uuid.New(),
			Name:         "create_order_saga",
			ReplyChannel: "kfk.dev.create_order_saga.reply",
			Steps:        steps,
		}
		message := Message{
			GlobalID: uuid.New(),
			EventID:  uuid.New(),
			EventType: EventType{
				SagaName: workflow.Name,
				StepName: createOrderStep.Name,
				Action:   SuccessActionType,
			},
			Saga: Saga{
				ID:           workflow.ID,
				Name:         workflow.Name,
				ReplyChannel: workflow.ReplyChannel,
				Step: SagaStep{
					ID:     createOrderStep.ID,
					Name:   createOrderStep.Name,
					Action: SuccessActionType,
				},
			},
			Metadata: map[string]string{
				"client_id": uuid.NewString(),
			},
			EventData: map[string]interface{}{},
		}

		step, err := workflow.GetNextStep(context.Background(), message)
		assert.Nil(t, step)
		assert.Nil(t, err)
	})

	t.Run("for successful message, should return next step when exists", func(t *testing.T) {
		steps := NewStepList()
		createOrderStep := steps.Append(&StepData{
			ID:          uuid.New(),
			Name:        "create_order",
			ServiceName: "order",
			Compensable: true,
			PayloadBuilder: func(ctx context.Context, data any) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
		})

		verifyClient := steps.Append(&StepData{
			ID:          uuid.New(),
			Name:        "verify_client",
			ServiceName: "client",
			Compensable: true,
			PayloadBuilder: func(ctx context.Context, data any) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
		})

		workflow := Workflow{
			ID:           uuid.New(),
			Name:         "create_order_saga",
			ReplyChannel: "kfk.dev.create_order_saga.reply",
			Steps:        steps,
		}

		message := Message{
			GlobalID: uuid.New(),
			EventID:  uuid.New(),
			EventType: EventType{
				SagaName: workflow.Name,
				StepName: createOrderStep.Name,
				Action:   SuccessActionType,
			},
			Saga: Saga{
				ID:           workflow.ID,
				Name:         workflow.Name,
				ReplyChannel: workflow.ReplyChannel,
				Step: SagaStep{
					ID:     createOrderStep.ID,
					Name:   createOrderStep.Name,
					Action: SuccessActionType,
				},
			},
			Metadata: map[string]string{
				"client_id": uuid.NewString(),
			},
			EventData: map[string]interface{}{},
		}

		step, err := workflow.GetNextStep(context.Background(), message)
		assert.Equal(t, verifyClient, step)
		assert.Nil(t, err)
	})

	t.Run("for error message type, should return the first compensable step [itself]", func(t *testing.T) {
		steps := NewStepList()
		createOrderStep := steps.Append(&StepData{
			ID:          uuid.New(),
			Name:        "create_order",
			ServiceName: "order",
			Compensable: true,
			PayloadBuilder: func(ctx context.Context, data any) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
		})

		_ = steps.Append(&StepData{
			ID:          uuid.New(),
			Name:        "verify_client",
			ServiceName: "client",
			Compensable: false,
			PayloadBuilder: func(ctx context.Context, data any) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
		})

		workflow := Workflow{
			ID:           uuid.New(),
			Name:         "create_order_saga",
			ReplyChannel: "kfk.dev.create_order_saga.reply",
			Steps:        steps,
		}

		message := Message{
			GlobalID: uuid.New(),
			EventID:  uuid.New(),
			EventType: EventType{
				SagaName: workflow.Name,
				StepName: createOrderStep.Name,
				Action:   FailureActionType,
			},
			Saga: Saga{
				ID:           workflow.ID,
				Name:         workflow.Name,
				ReplyChannel: workflow.ReplyChannel,
				Step: SagaStep{
					ID:     createOrderStep.ID,
					Name:   createOrderStep.Name,
					Action: FailureActionType,
				},
			},
			Metadata: map[string]string{
				"client_id": uuid.NewString(),
			},
			EventData: map[string]interface{}{},
		}

		step, err := workflow.GetNextStep(context.Background(), message)
		assert.Equal(t, createOrderStep, step)
		assert.Nil(t, err)
	})

	t.Run("for error message type, should return the the first compensable step [previous]", func(t *testing.T) {
		steps := NewStepList()
		createOrderStep := steps.Append(&StepData{
			ID:          uuid.New(),
			Name:        "create_order",
			ServiceName: "order",
			Compensable: true,
			PayloadBuilder: func(ctx context.Context, data any) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
		})

		verifyStep := steps.Append(&StepData{
			ID:          uuid.New(),
			Name:        "verify_client",
			ServiceName: "client",
			Compensable: false,
			PayloadBuilder: func(ctx context.Context, data any) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
		})

		workflow := Workflow{
			ID:           uuid.New(),
			Name:         "create_order_saga",
			ReplyChannel: "kfk.dev.create_order_saga.reply",
			Steps:        steps,
		}

		message := Message{
			GlobalID: uuid.New(),
			EventID:  uuid.New(),
			EventType: EventType{
				SagaName: workflow.Name,
				StepName: verifyStep.Name,
				Action:   FailureActionType,
			},
			Saga: Saga{
				ID:           workflow.ID,
				Name:         workflow.Name,
				ReplyChannel: workflow.ReplyChannel,
				Step: SagaStep{
					ID:     verifyStep.ID,
					Name:   verifyStep.Name,
					Action: FailureActionType,
				},
			},
			Metadata: map[string]string{
				"client_id": uuid.NewString(),
			},
			EventData: map[string]interface{}{},
		}

		step, err := workflow.GetNextStep(context.Background(), message)
		assert.Equal(t, createOrderStep, step)
		assert.Nil(t, err)
	})

	t.Run("for error message type, should return nil if there are no compensable steps", func(t *testing.T) {
		steps := NewStepList()
		_ = steps.Append(&StepData{
			ID:          uuid.New(),
			Name:        "create_order",
			ServiceName: "order",
			Compensable: false,
			PayloadBuilder: func(ctx context.Context, data any) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
		})

		verifyStep := steps.Append(&StepData{
			ID:          uuid.New(),
			Name:        "verify_client",
			ServiceName: "client",
			Compensable: false,
			PayloadBuilder: func(ctx context.Context, data any) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
		})

		workflow := Workflow{
			ID:           uuid.New(),
			Name:         "create_order_saga",
			ReplyChannel: "kfk.dev.create_order_saga.reply",
			Steps:        steps,
		}

		message := Message{
			GlobalID: uuid.New(),
			EventID:  uuid.New(),
			EventType: EventType{
				SagaName: workflow.Name,
				StepName: verifyStep.Name,
				Action:   FailureActionType,
			},
			Saga: Saga{
				ID:           workflow.ID,
				Name:         workflow.Name,
				ReplyChannel: workflow.ReplyChannel,
				Step: SagaStep{
					ID:     verifyStep.ID,
					Name:   verifyStep.Name,
					Action: FailureActionType,
				},
			},
			Metadata: map[string]string{
				"client_id": uuid.NewString(),
			},
			EventData: map[string]interface{}{},
		}

		step, err := workflow.GetNextStep(context.Background(), message)
		assert.Nil(t, step)
		assert.Nil(t, err)
	})

	t.Run("for compensated message type, should return nil if there are no more compensable steps", func(t *testing.T) {
		steps := NewStepList()
		_ = steps.Append(&StepData{
			ID:          uuid.New(),
			Name:        "create_order",
			ServiceName: "order",
			Compensable: false,
			PayloadBuilder: func(ctx context.Context, data any) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
		})

		verifyStep := steps.Append(&StepData{
			ID:          uuid.New(),
			Name:        "verify_client",
			ServiceName: "client",
			Compensable: false,
			PayloadBuilder: func(ctx context.Context, data any) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
		})

		workflow := Workflow{
			ID:           uuid.New(),
			Name:         "create_order_saga",
			ReplyChannel: "kfk.dev.create_order_saga.reply",
			Steps:        steps,
		}

		message := Message{
			GlobalID: uuid.New(),
			EventID:  uuid.New(),
			EventType: EventType{
				SagaName: workflow.Name,
				StepName: verifyStep.Name,
				Action:   CompensatedActionType,
			},
			Saga: Saga{
				ID:           workflow.ID,
				Name:         workflow.Name,
				ReplyChannel: workflow.ReplyChannel,
				Step: SagaStep{
					ID:     verifyStep.ID,
					Name:   verifyStep.Name,
					Action: CompensatedActionType,
				},
			},
			Metadata: map[string]string{
				"client_id": uuid.NewString(),
			},
			EventData: map[string]interface{}{},
		}

		step, err := workflow.GetNextStep(context.Background(), message)
		assert.Nil(t, step)
		assert.Nil(t, err)
	})

	t.Run("for compensated message type, should return the first compensable step", func(t *testing.T) {
		steps := NewStepList()
		createOrderStep := steps.Append(&StepData{
			ID:          uuid.New(),
			Name:        "create_order",
			ServiceName: "order",
			Compensable: true,
			PayloadBuilder: func(ctx context.Context, data any) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
		})

		verifyStep := steps.Append(&StepData{
			ID:          uuid.New(),
			Name:        "verify_client",
			ServiceName: "client",
			Compensable: false,
			PayloadBuilder: func(ctx context.Context, data any) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
		})

		workflow := Workflow{
			ID:           uuid.New(),
			Name:         "create_order_saga",
			ReplyChannel: "kfk.dev.create_order_saga.reply",
			Steps:        steps,
		}

		message := Message{
			GlobalID: uuid.New(),
			EventID:  uuid.New(),
			EventType: EventType{
				SagaName: workflow.Name,
				StepName: verifyStep.Name,
				Action:   CompensatedActionType,
			},
			Saga: Saga{
				ID:           workflow.ID,
				Name:         workflow.Name,
				ReplyChannel: workflow.ReplyChannel,
				Step: SagaStep{
					ID:     verifyStep.ID,
					Name:   verifyStep.Name,
					Action: CompensatedActionType,
				},
			},
			Metadata: map[string]string{
				"client_id": uuid.NewString(),
			},
			EventData: map[string]interface{}{},
		}

		step, err := workflow.GetNextStep(context.Background(), message)
		assert.Equal(t, createOrderStep, step)
		assert.Nil(t, err)
	})

	t.Run("for ErrUnknownActionType if event action type is not processable", func(t *testing.T) {
		steps := NewStepList()
		createOrderStep := steps.Append(&StepData{
			ID:          uuid.New(),
			Name:        "create_order",
			ServiceName: "order",
			Compensable: true,
			PayloadBuilder: func(ctx context.Context, data any) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
		})

		workflow := Workflow{
			ID:           uuid.New(),
			Name:         "create_order_saga",
			ReplyChannel: "kfk.dev.create_order_saga.reply",
			Steps:        steps,
		}

		message := Message{
			GlobalID: uuid.New(),
			EventID:  uuid.New(),
			EventType: EventType{
				SagaName: workflow.Name,
				StepName: createOrderStep.Name,
				Action:   RequestActionType,
			},
			Saga: Saga{
				ID:           workflow.ID,
				Name:         workflow.Name,
				ReplyChannel: workflow.ReplyChannel,
				Step: SagaStep{
					ID:     createOrderStep.ID,
					Name:   createOrderStep.Name,
					Action: RequestActionType,
				},
			},
			Metadata: map[string]string{
				"client_id": uuid.NewString(),
			},
			EventData: map[string]interface{}{},
		}

		step, err := workflow.GetNextStep(context.Background(), message)
		assert.Nil(t, step)
		assert.Error(t, err)
		assert.Equal(t, ErrUnknownActionType, err)
	})
}

func TestWorkflow_ConsumerEventTypes(t *testing.T) {
	t.Run("should return a empty map if there are no steps in the workflow", func(t *testing.T) {
		workflow := Workflow{
			ID:           uuid.New(),
			Name:         "create_order_saga",
			ReplyChannel: "kfk.dev.create_order_saga.reply",
			Steps:        NewStepList(),
		}

		eventTypes := workflow.ConsumerEventTypes()
		assert.Empty(t, eventTypes)
	})
	t.Run("should return map of event types for consumer", func(t *testing.T) {
		steps := NewStepList(
			&StepData{
				ID:          uuid.New(),
				Name:        "create_order",
				ServiceName: "order",
				Compensable: true,
			},
			&StepData{
				ID:          uuid.New(),
				Name:        "verify_client",
				ServiceName: "client",
				Compensable: true,
			},
		)

		workflow := Workflow{
			ID:           uuid.New(),
			Name:         "create_order_saga",
			ReplyChannel: "kfk.dev.create_order_saga.reply",
			Steps:        steps,
		}

		eventTypes := workflow.ConsumerEventTypes()
		assert.Equal(t, len(eventTypes), 6)
		for _, step := range eventTypes {
			assert.Equal(t, workflow.ID, step)
		}
		expectedEventTypes := []string{
			"create_order_saga.create_order.success",
			"create_order_saga.create_order.failure",
			"create_order_saga.create_order.compensated",
			"create_order_saga.verify_client.success",
			"create_order_saga.verify_client.failure",
			"create_order_saga.verify_client.compensated",
		}
		for _, eventType := range expectedEventTypes {
			_, ok := eventTypes[eventType]
			assert.True(t, ok)
		}

	})
}
