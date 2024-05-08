package handlers

import (
	"context"
	"encoding/json"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application/usecases"
	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CreateOrderHandler struct {
	logger             *zap.SugaredLogger
	createOrderUseCase usecases.CreateOrderUseCasePort
}

var (
	_ application.MessageHandler = (*CreateOrderHandler)(nil)
)

type request struct {
	CustomerID   uuid.UUID `json:"customer_id"`
	Amount       *int64    `json:"amount"`
	CurrencyCode string    `json:"currency_code"`
}

func NewCreateOrderHandler(logger *zap.SugaredLogger, createOrderUseCase usecases.CreateOrderUseCasePort) *CreateOrderHandler {
	return &CreateOrderHandler{
		logger:             logger,
		createOrderUseCase: createOrderUseCase,
	}
}

func parseInput(data map[string]interface{}, dest interface{}) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw, dest)
	if err != nil {
		return err
	}
	return nil
}

func (h *CreateOrderHandler) Handle(ctx context.Context, msg *saga.Message) (*saga.Message, error) {
	lggr := h.logger
	lggr.Infof("Handling message [%s]", msg.EventType.String())

	globalID := msg.GlobalID
	var req request
	err := parseInput(msg.EventData, &req)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error reading input")
		return nil, err
	}
	createRes, err := h.createOrderUseCase.Execute(ctx, usecases.CreateOrderRequest{
		GlobalID:     globalID,
		CustomerID:   req.CustomerID,
		Amount:       *req.Amount,
		CurrencyCode: req.CurrencyCode,
	})
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error creating order")
		resEventType := saga.NewReplyEventType(msg.EventType, saga.FailureActionType)
		replyMessage := saga.NewParticipantMessage(globalID, nil, nil, resEventType, msg.Saga.ReplyChannel)
		return replyMessage, nil
	}

	lggr.Infof("Successfully created order [%s]", createRes)
	res := map[string]interface{}{"order_id": createRes.ID}
	resEventType := saga.NewReplyEventType(msg.EventType, saga.SuccessActionType)
	replyMessage := saga.NewParticipantMessage(globalID, res, nil, resEventType, msg.Saga.ReplyChannel)
	return replyMessage, nil
}
