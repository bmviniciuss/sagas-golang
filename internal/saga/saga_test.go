package saga

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWorkflow_GetNextStep(t *testing.T) {
	t.Run("should return error when current step is not found in message workflow", func(t *testing.T) {
		steps := NewStepList()
		createOrderStep := &StepData{
			ID:          uuid.New(),
			Name:        "create_order",
			ServiceName: "order",
			Compensable: true,
			PayloadBuilder: func(ctx context.Context, data any) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
		}

		steps.Append(createOrderStep)
		workflow := Workflow{
			ID:           uuid.New(),
			Name:         "create_order_saga",
			ReplyChannel: "kfk.dev.create_order_saga.reply",
			Steps:        steps,
		}

		message := Message{
			EventID:   uuid.NewString(),
			EventType: "create_order.success",
			Metadata: Metadata{
				GlobalID: uuid.NewString(),
				ClientID: uuid.NewString(),
			},
			EventData: map[string]interface{}{},
			Saga: Saga{
				ID:           workflow.ID,
				Name:         workflow.Name,
				ReplyChannel: workflow.ReplyChannel,
				Step: SagaStep{
					ID:   uuid.New(),
					Name: "any-other-step",
				},
			},
			ActionType: SuccessActionType,
		}

		step, err := workflow.GetNextStep(context.Background(), message)
		assert.Nil(t, step)
		assert.Error(t, err)
		assert.Equal(t, ErrCurrentStepNotFound, err)
	})

	t.Run("should return nil when message action type is success and there are no more steps", func(t *testing.T) {
		steps := NewStepList()
		createOrderStep := &StepData{
			ID:          uuid.New(),
			Name:        "create_order",
			ServiceName: "order",
			Compensable: true,
			PayloadBuilder: func(ctx context.Context, data any) (map[string]interface{}, error) {
				return map[string]interface{}{}, nil
			},
		}
		steps.Append(createOrderStep)

		workflow := Workflow{
			ID:           uuid.New(),
			Name:         "create_order_saga",
			ReplyChannel: "kfk.dev.create_order_saga.reply",
			Steps:        steps,
		}

		message := Message{
			EventID:   uuid.NewString(),
			EventType: "create_order.success",
			Metadata: Metadata{
				GlobalID: uuid.NewString(),
				ClientID: uuid.NewString(),
			},
			EventData: map[string]interface{}{},
			Saga: Saga{
				ID:           workflow.ID,
				Name:         workflow.Name,
				ReplyChannel: workflow.ReplyChannel,
				Step: SagaStep{
					ID:   createOrderStep.ID,
					Name: createOrderStep.Name,
				},
			},
			ActionType: SuccessActionType,
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
			EventID:   uuid.NewString(),
			EventType: "create_order.success",
			Metadata: Metadata{
				GlobalID: uuid.NewString(),
				ClientID: uuid.NewString(),
			},
			EventData: map[string]interface{}{},
			Saga: Saga{
				ID:           workflow.ID,
				Name:         workflow.Name,
				ReplyChannel: workflow.ReplyChannel,
				Step: SagaStep{
					ID:   createOrderStep.ID,
					Name: createOrderStep.Name,
				},
			},
			ActionType: SuccessActionType,
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
			EventID:   uuid.NewString(),
			EventType: "create_order.success",
			Metadata: Metadata{
				GlobalID: uuid.NewString(),
				ClientID: uuid.NewString(),
			},
			EventData: map[string]interface{}{},
			Saga: Saga{
				ID:           workflow.ID,
				Name:         workflow.Name,
				ReplyChannel: workflow.ReplyChannel,
				Step: SagaStep{
					ID:   createOrderStep.ID,
					Name: createOrderStep.Name,
				},
			},
			ActionType: FailureActionType,
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
			EventID:   uuid.NewString(),
			EventType: "verify_client.failure",
			Metadata: Metadata{
				GlobalID: uuid.NewString(),
				ClientID: uuid.NewString(),
			},
			EventData: map[string]interface{}{},
			Saga: Saga{
				ID:           workflow.ID,
				Name:         workflow.Name,
				ReplyChannel: workflow.ReplyChannel,
				Step: SagaStep{
					ID:   verifyStep.ID,
					Name: verifyStep.Name,
				},
			},
			ActionType: FailureActionType,
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
			EventID:   uuid.NewString(),
			EventType: "verify_client.success",
			Metadata: Metadata{
				GlobalID: uuid.NewString(),
				ClientID: uuid.NewString(),
			},
			EventData: map[string]interface{}{},
			Saga: Saga{
				ID:           workflow.ID,
				Name:         workflow.Name,
				ReplyChannel: workflow.ReplyChannel,
				Step: SagaStep{
					ID:   verifyStep.ID,
					Name: verifyStep.Name,
				},
			},
			ActionType: FailureActionType,
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
			EventID:   uuid.NewString(),
			EventType: "verify_client.success",
			Metadata: Metadata{
				GlobalID: uuid.NewString(),
				ClientID: uuid.NewString(),
			},
			EventData: map[string]interface{}{},
			Saga: Saga{
				ID:           workflow.ID,
				Name:         workflow.Name,
				ReplyChannel: workflow.ReplyChannel,
				Step: SagaStep{
					ID:   verifyStep.ID,
					Name: verifyStep.Name,
				},
			},
			ActionType: CompensatedActionType,
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
			EventID:   uuid.NewString(),
			EventType: "verify_client.success",
			Metadata: Metadata{
				GlobalID: uuid.NewString(),
				ClientID: uuid.NewString(),
			},
			EventData: map[string]interface{}{},
			Saga: Saga{
				ID:           workflow.ID,
				Name:         workflow.Name,
				ReplyChannel: workflow.ReplyChannel,
				Step: SagaStep{
					ID:   verifyStep.ID,
					Name: verifyStep.Name,
				},
			},
			ActionType: CompensatedActionType,
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
			EventID:   uuid.NewString(),
			EventType: "create_order.success",
			Metadata: Metadata{
				GlobalID: uuid.NewString(),
				ClientID: uuid.NewString(),
			},
			EventData: map[string]interface{}{},
			Saga: Saga{
				ID:           workflow.ID,
				Name:         workflow.Name,
				ReplyChannel: workflow.ReplyChannel,
				Step: SagaStep{
					ID:   createOrderStep.ID,
					Name: createOrderStep.Name,
				},
			},
			ActionType: RequestActionType,
		}

		step, err := workflow.GetNextStep(context.Background(), message)
		assert.Nil(t, step)
		assert.Error(t, err)
		assert.Equal(t, ErrUnknownActionType, err)
	})
}
