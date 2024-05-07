package saga

// import (
// 	"context"
// 	"testing"

// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// )

// type payloadBuilderMock struct{}

// func (pb *payloadBuilderMock) Build(ctx context.Context, data map[string]interface{}, action ActionType) (map[string]interface{}, error) {
// 	return map[string]interface{}{}, nil
// }

// func TestWorkflow_GetNextStep(t *testing.T) {

// 	t.Run("should return error when current step is not found in message workflow", func(t *testing.T) {
// 		steps := NewStepList(&StepData{
// 			Name:           "create_order",
// 			ServiceName:    "order",
// 			Compensable:    true,
// 			PayloadBuilder: &payloadBuilderMock{},
// 		})
// 		workflow := Workflow{
// 			Name:         "create_order_saga",
// 			ReplyChannel: "kfk.dev.create_order_saga.reply",
// 			Steps:        steps,
// 		}
// 		message := Message{
// 			GlobalID: uuid.New(),
// 			EventID:  uuid.New(),
// 			EventType: EventType{
// 				SagaName: workflow.Name,
// 				StepName: "create-order",
// 				Action:   SuccessActionType,
// 			},
// 			Saga: Saga{
// 				Name:         workflow.Name,
// 				ReplyChannel: workflow.ReplyChannel,
// 				Step: SagaStep{
// 					Name:   "create-order",
// 					Action: SuccessActionType,
// 				},
// 			},
// 			Metadata: map[string]string{
// 				"client_id": uuid.NewString(),
// 			},
// 			EventData: map[string]interface{}{},
// 		}
// 		step, err := workflow.GetNextStep(context.Background(), message)
// 		assert.Nil(t, step)
// 		assert.Error(t, err)
// 		assert.Equal(t, ErrCurrentStepNotFound, err)
// 	})

// 	t.Run("should return nil when message action type is success and there are no more steps", func(t *testing.T) {
// 		createOrderStep := &StepData{
// 			Name:           "create_order",
// 			ServiceName:    "order",
// 			Compensable:    true,
// 			PayloadBuilder: &payloadBuilderMock{},
// 		}
// 		steps := NewStepList(createOrderStep)
// 		workflow := Workflow{
// 			Name:         "create_order_saga",
// 			ReplyChannel: "kfk.dev.create_order_saga.reply",
// 			Steps:        steps,
// 		}
// 		message := Message{
// 			GlobalID: uuid.New(),
// 			EventID:  uuid.New(),
// 			EventType: EventType{
// 				SagaName: workflow.Name,
// 				StepName: createOrderStep.Name,
// 				Action:   SuccessActionType,
// 			},
// 			Saga: Saga{
// 				Name:         workflow.Name,
// 				ReplyChannel: workflow.ReplyChannel,
// 				Step: SagaStep{
// 					Name:   createOrderStep.Name,
// 					Action: SuccessActionType,
// 				},
// 			},
// 			Metadata: map[string]string{
// 				"client_id": uuid.NewString(),
// 			},
// 			EventData: map[string]interface{}{},
// 		}

// 		step, err := workflow.GetNextStep(context.Background(), message)
// 		assert.Nil(t, step)
// 		assert.Nil(t, err)
// 	})

// 	t.Run("for successful message, should return next step when exists", func(t *testing.T) {
// 		steps := NewStepList()
// 		createOrderStep := steps.Append(&StepData{
// 			Name:           "create_order",
// 			ServiceName:    "order",
// 			Compensable:    true,
// 			PayloadBuilder: &payloadBuilderMock{},
// 		})

// 		verifyClient := steps.Append(&StepData{
// 			Name:           "verify_client",
// 			ServiceName:    "client",
// 			Compensable:    true,
// 			PayloadBuilder: &payloadBuilderMock{},
// 		})

// 		workflow := Workflow{
// 			Name:         "create_order_saga",
// 			ReplyChannel: "kfk.dev.create_order_saga.reply",
// 			Steps:        steps,
// 		}

// 		message := Message{
// 			GlobalID: uuid.New(),
// 			EventID:  uuid.New(),
// 			EventType: EventType{
// 				SagaName: workflow.Name,
// 				StepName: createOrderStep.Name,
// 				Action:   SuccessActionType,
// 			},
// 			Saga: Saga{
// 				Name:         workflow.Name,
// 				ReplyChannel: workflow.ReplyChannel,
// 				Step: SagaStep{
// 					Name:   createOrderStep.Name,
// 					Action: SuccessActionType,
// 				},
// 			},
// 			Metadata: map[string]string{
// 				"client_id": uuid.NewString(),
// 			},
// 			EventData: map[string]interface{}{},
// 		}

// 		step, err := workflow.GetNextStep(context.Background(), message)
// 		assert.Equal(t, verifyClient, step)
// 		assert.Nil(t, err)
// 	})

// 	t.Run("for error message type, should return the first compensable step [itself]", func(t *testing.T) {
// 		steps := NewStepList()
// 		createOrderStep := steps.Append(&StepData{
// 			Name:           "create_order",
// 			ServiceName:    "order",
// 			Compensable:    true,
// 			PayloadBuilder: &payloadBuilderMock{},
// 		})

// 		_ = steps.Append(&StepData{
// 			Name:           "verify_client",
// 			ServiceName:    "client",
// 			Compensable:    false,
// 			PayloadBuilder: &payloadBuilderMock{},
// 		})

// 		workflow := Workflow{
// 			Name:         "create_order_saga",
// 			ReplyChannel: "kfk.dev.create_order_saga.reply",
// 			Steps:        steps,
// 		}

// 		message := Message{
// 			GlobalID: uuid.New(),
// 			EventID:  uuid.New(),
// 			EventType: EventType{
// 				SagaName: workflow.Name,
// 				StepName: createOrderStep.Name,
// 				Action:   FailureActionType,
// 			},
// 			Saga: Saga{
// 				Name:         workflow.Name,
// 				ReplyChannel: workflow.ReplyChannel,
// 				Step: SagaStep{
// 					Name:   createOrderStep.Name,
// 					Action: FailureActionType,
// 				},
// 			},
// 			Metadata: map[string]string{
// 				"client_id": uuid.NewString(),
// 			},
// 			EventData: map[string]interface{}{},
// 		}

// 		step, err := workflow.GetNextStep(context.Background(), message)
// 		assert.Equal(t, createOrderStep, step)
// 		assert.Nil(t, err)
// 	})

// 	t.Run("for error message type, should return the the first compensable step [previous]", func(t *testing.T) {
// 		steps := NewStepList()
// 		createOrderStep := steps.Append(&StepData{
// 			Name:           "create_order",
// 			ServiceName:    "order",
// 			Compensable:    true,
// 			PayloadBuilder: &payloadBuilderMock{},
// 		})

// 		verifyStep := steps.Append(&StepData{
// 			Name:           "verify_client",
// 			ServiceName:    "client",
// 			Compensable:    false,
// 			PayloadBuilder: &payloadBuilderMock{},
// 		})

// 		workflow := Workflow{
// 			Name:         "create_order_saga",
// 			ReplyChannel: "kfk.dev.create_order_saga.reply",
// 			Steps:        steps,
// 		}

// 		message := Message{
// 			GlobalID: uuid.New(),
// 			EventID:  uuid.New(),
// 			EventType: EventType{
// 				SagaName: workflow.Name,
// 				StepName: verifyStep.Name,
// 				Action:   FailureActionType,
// 			},
// 			Saga: Saga{
// 				Name:         workflow.Name,
// 				ReplyChannel: workflow.ReplyChannel,
// 				Step: SagaStep{
// 					Name:   verifyStep.Name,
// 					Action: FailureActionType,
// 				},
// 			},
// 			Metadata: map[string]string{
// 				"client_id": uuid.NewString(),
// 			},
// 			EventData: map[string]interface{}{},
// 		}

// 		step, err := workflow.GetNextStep(context.Background(), message)
// 		assert.Equal(t, createOrderStep, step)
// 		assert.Nil(t, err)
// 	})

// 	t.Run("for error message type, should return nil if there are no compensable steps", func(t *testing.T) {
// 		steps := NewStepList()
// 		_ = steps.Append(&StepData{
// 			Name:           "create_order",
// 			ServiceName:    "order",
// 			Compensable:    false,
// 			PayloadBuilder: &payloadBuilderMock{},
// 		})

// 		verifyStep := steps.Append(&StepData{
// 			Name:           "verify_client",
// 			ServiceName:    "client",
// 			Compensable:    false,
// 			PayloadBuilder: &payloadBuilderMock{},
// 		})

// 		workflow := Workflow{
// 			Name:         "create_order_saga",
// 			ReplyChannel: "kfk.dev.create_order_saga.reply",
// 			Steps:        steps,
// 		}

// 		message := Message{
// 			GlobalID: uuid.New(),
// 			EventID:  uuid.New(),
// 			EventType: EventType{
// 				SagaName: workflow.Name,
// 				StepName: verifyStep.Name,
// 				Action:   FailureActionType,
// 			},
// 			Saga: Saga{
// 				Name:         workflow.Name,
// 				ReplyChannel: workflow.ReplyChannel,
// 				Step: SagaStep{
// 					Name:   verifyStep.Name,
// 					Action: FailureActionType,
// 				},
// 			},
// 			Metadata: map[string]string{
// 				"client_id": uuid.NewString(),
// 			},
// 			EventData: map[string]interface{}{},
// 		}

// 		step, err := workflow.GetNextStep(context.Background(), message)
// 		assert.Nil(t, step)
// 		assert.Nil(t, err)
// 	})

// 	t.Run("for compensated message type, should return nil if there are no more compensable steps", func(t *testing.T) {
// 		steps := NewStepList()
// 		_ = steps.Append(&StepData{
// 			Name:           "create_order",
// 			ServiceName:    "order",
// 			Compensable:    false,
// 			PayloadBuilder: &payloadBuilderMock{},
// 		})

// 		verifyStep := steps.Append(&StepData{
// 			Name:           "verify_client",
// 			ServiceName:    "client",
// 			Compensable:    false,
// 			PayloadBuilder: &payloadBuilderMock{},
// 		})

// 		workflow := Workflow{
// 			Name:         "create_order_saga",
// 			ReplyChannel: "kfk.dev.create_order_saga.reply",
// 			Steps:        steps,
// 		}

// 		message := Message{
// 			GlobalID: uuid.New(),
// 			EventID:  uuid.New(),
// 			EventType: EventType{
// 				SagaName: workflow.Name,
// 				StepName: verifyStep.Name,
// 				Action:   CompensatedActionType,
// 			},
// 			Saga: Saga{
// 				Name:         workflow.Name,
// 				ReplyChannel: workflow.ReplyChannel,
// 				Step: SagaStep{
// 					Name:   verifyStep.Name,
// 					Action: CompensatedActionType,
// 				},
// 			},
// 			Metadata: map[string]string{
// 				"client_id": uuid.NewString(),
// 			},
// 			EventData: map[string]interface{}{},
// 		}

// 		step, err := workflow.GetNextStep(context.Background(), message)
// 		assert.Nil(t, step)
// 		assert.Nil(t, err)
// 	})

// 	t.Run("for compensated message type, should return the first compensable step", func(t *testing.T) {
// 		steps := NewStepList()
// 		createOrderStep := steps.Append(&StepData{
// 			Name:           "create_order",
// 			ServiceName:    "order",
// 			Compensable:    true,
// 			PayloadBuilder: &payloadBuilderMock{},
// 		})

// 		verifyStep := steps.Append(&StepData{
// 			Name:           "verify_client",
// 			ServiceName:    "client",
// 			Compensable:    false,
// 			PayloadBuilder: &payloadBuilderMock{},
// 		})

// 		workflow := Workflow{
// 			Name:         "create_order_saga",
// 			ReplyChannel: "kfk.dev.create_order_saga.reply",
// 			Steps:        steps,
// 		}

// 		message := Message{
// 			GlobalID: uuid.New(),
// 			EventID:  uuid.New(),
// 			EventType: EventType{
// 				SagaName: workflow.Name,
// 				StepName: verifyStep.Name,
// 				Action:   CompensatedActionType,
// 			},
// 			Saga: Saga{
// 				Name:         workflow.Name,
// 				ReplyChannel: workflow.ReplyChannel,
// 				Step: SagaStep{
// 					Name:   verifyStep.Name,
// 					Action: CompensatedActionType,
// 				},
// 			},
// 			Metadata: map[string]string{
// 				"client_id": uuid.NewString(),
// 			},
// 			EventData: map[string]interface{}{},
// 		}

// 		step, err := workflow.GetNextStep(context.Background(), message)
// 		assert.Equal(t, createOrderStep, step)
// 		assert.Nil(t, err)
// 	})

// 	t.Run("for ErrUnknownActionType if event action type is not processable", func(t *testing.T) {
// 		steps := NewStepList()
// 		createOrderStep := steps.Append(&StepData{
// 			Name:           "create_order",
// 			ServiceName:    "order",
// 			Compensable:    true,
// 			PayloadBuilder: &payloadBuilderMock{},
// 		})

// 		workflow := Workflow{
// 			Name:         "create_order_saga",
// 			ReplyChannel: "kfk.dev.create_order_saga.reply",
// 			Steps:        steps,
// 		}

// 		message := Message{
// 			GlobalID: uuid.New(),
// 			EventID:  uuid.New(),
// 			EventType: EventType{
// 				SagaName: workflow.Name,
// 				StepName: createOrderStep.Name,
// 				Action:   RequestActionType,
// 			},
// 			Saga: Saga{
// 				Name:         workflow.Name,
// 				ReplyChannel: workflow.ReplyChannel,
// 				Step: SagaStep{
// 					Name:   createOrderStep.Name,
// 					Action: RequestActionType,
// 				},
// 			},
// 			Metadata: map[string]string{
// 				"client_id": uuid.NewString(),
// 			},
// 			EventData: map[string]interface{}{},
// 		}

// 		step, err := workflow.GetNextStep(context.Background(), message)
// 		assert.Nil(t, step)
// 		assert.Error(t, err)
// 		assert.Equal(t, ErrUnknownActionType, err)
// 	})
// }
