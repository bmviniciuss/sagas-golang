package main

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"

// 	"github.com/bmviniciuss/sagas-golang/internal/saga"
// 	"github.com/bmviniciuss/sagas-golang/internal/streaming"
// 	"github.com/google/uuid"
// 	"github.com/go-viper/mapstructure/v2"
// )

// type CreateOrderInput struct {
// 	Amount int64 `mapstructure:"amount"`
// }

// func main() {
// 	reqBody := map[string]interface{}{
// 		"amount": int64(1000),
// 	}

// 	stepsList := saga.NewStepList()
// 	stepsList.Append(&saga.StepData{
// 		ID:          uuid.New(),
// 		Name:        "create_order",
// 		ServiceName: "order",
// 		Compensable: true,
// 		PayloadBuilder: func(ctx context.Context, data any) (map[string]interface{}, error) {
// 			var input CreateOrderInput
// 			err := mapstructure.Decode(data, &input)
// 			if err != nil {
// 				return nil, err
// 			}
// 			return map[string]interface{}{
// 				"amount": input.Amount,
// 			}, nil
// 		},
// 	})

// 	createOrderSaga := saga.Workflow{
// 		ID:           uuid.New(),
// 		Name:         "create_order_saga",
// 		ReplyChannel: "kfk.dev.create_order_saga.reply",
// 		Steps:        stepsList,
// 	}

// 	StartSaga(context.TODO(), createOrderSaga, reqBody)
// 	consumer := Consumer{}
// 	ws := map[uuid.UUID]saga.Workflow{
// 		createOrderSaga.ID: createOrderSaga,
// 	}
// 	err := consumer.Start(context.Background(), ws)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func StartSaga(ctx context.Context, saga saga.Workflow, data any) error {
// 	fmt.Println("Starting Saga")
// 	firstStep, ok := saga.Steps.Head()
// 	if !ok {
// 		return errors.New("saga does not contains first step")
// 	}

// 	globalID := uuid.New()
// 	eventID := uuid.New()
// 	eventData, err := firstStep.PayloadBuilder(ctx, data)
// 	if err != nil {
// 		fmt.Println("Error building first step payload")
// 		return err
// 	}

// 	msg := streaming.Message{
// 		Metadata: streaming.Metadata{
// 			GlobalID: globalID.String(),
// 			ClientID: uuid.NewString(),
// 		},
// 		EventID:   eventID.String(),
// 		EventType: fmt.Sprintf("%s.%s", saga.Name, firstStep.Name),
// 		EventData: eventData,
// 		Saga: streaming.Saga{
// 			ID:           saga.ID,
// 			Name:         saga.Name,
// 			ReplyChannel: saga.ReplyChannel,
// 			Step: streaming.SagaStep{
// 				ID:   firstStep.ID,
// 				Name: firstStep.Name,
// 			},
// 			Metadata: map[string]interface{}{},
// 		},
// 	}

// 	msgBytes, _ := json.Marshal(msg)

// 	fmt.Printf("Built First Step event data = [%s]\n\n", string(msgBytes))
// 	return nil
// }

// type Consumer struct {
// }

// func pool() []byte {
// 	return []byte(`{
// 		"metadata": {
// 			"global_id": "2780cf76-1829-42b2-bdb6-40c41b6a9ccb",
// 			"client_id": "7c1b2ca5-cfde-4ea7-adb7-9a06d34125f9"
// 		},
// 		"event_id": "f1763d99-4bdd-41be-acee-c71e8d1c8d22",
// 		"event_type": "create_order_saga.create_order.success",
// 		"event_data": {
// 			"id": "a5475997-2458-41d1-a624-0a1e471bc0b9"
// 		},
// 		"saga": {
// 			"id": "22fb74c3-82e8-4140-ab99-d9c086581f2d",
// 			"name": "create_order_saga",
// 			"reply_channel": "kfk.dev.create_order_saga.reply",
// 			"step": {
// 				"id": "d8073562-d544-4290-aa00-4c919a554e27",
// 				"name": "create_order"
// 			},
// 			"metadata": {
// 				"verify_step": {
// 					"validation_external_id": "b7dd7bbe-6463-4b71-b33f-02682b9293dd"
// 				}
// 			}
// 		}
// 	}
// 	`)
// }

// func (c *Consumer) Start(ctx context.Context, workflows map[uuid.UUID]saga.Workflow) error {
// 	msg := pool()
// 	fmt.Println("Received message: ", string(msg))
// 	eventTypes := make(map[string]uuid.UUID)
// 	for _, workflow := range workflows {
// 		we := workflow.EventTypes()
// 		for k, v := range we {
// 			eventTypes[k] = v
// 		}
// 	}

// 	var message streaming.Message
// 	err := json.Unmarshal(msg, &message)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Printf("Parsed message [%+v]\n", message)

// 	actionType, ok := message.GetActionType()
// 	if !ok {
// 		return errors.New("unreconizeable event type action type")
// 	}

// 	workflowID, ok := eventTypes[message.EventType]
// 	if !ok {
// 		fmt.Println("Event should not be processed")
// 		return nil
// 	}
// 	fmt.Printf("Event with global_id [%s] with event_type of [%s] and id [%s] should be processed with workflow [%s]\n",
// 		message.Metadata.GlobalID,
// 		message.EventType,
// 		message.EventID,
// 		workflowID.String(),
// 	)

// 	workflow, ok := workflows[workflowID]
// 	if !ok {
// 		return errors.New("workflow not found")
// 	}
// 	fmt.Println("Workflow found")
// 	currentStep, ok := workflow.Steps.GetStep(message.Saga.Step.ID)
// 	if !ok {
// 		fmt.Println("Error getting saga current step")
// 		return errors.New("unreconized current step")
// 	}
// 	fmt.Println("Current step found")

// 	if actionType.IsSuccess() {
// 		fmt.Println("Action Type is success. Procedding to saga's next step")
// 		// next step
// 		_, ok := currentStep.Next()
// 		if !ok {
// 			fmt.Println("Saga does not have next step and has finished")
// 		}
// 		fmt.Println("Saga has next step and will process next step.")
// 	}
// 	if actionType.IsFailure() {
// 		fmt.Println("Action Type is failure. Start compensation flow from previous step")
// 		_, ok := currentStep.Previous() // TODO: get prev steps with compensation

// 	}

// 	// nextStep, ok := currentStep.Next()
// 	// if !ok {
// 	// 	fmt.Println("Saga has finished")
// 	// }
// 	// fmt.Println("Saga has next step")

// 	return nil
// }
