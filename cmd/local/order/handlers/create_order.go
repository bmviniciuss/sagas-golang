package handlers

import (
	"context"
	"encoding/json"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application/usecases"
	"github.com/bmviniciuss/sagas-golang/pkg/events"
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

// TODO: add enum
// Request:      "create_order",
// Success:      "order_created",
// Failure:      "order_creation_failed",
// Compensation: "order_creation_compensated",

func (h *CreateOrderHandler) Handle(ctx context.Context, event *events.Event) (*events.Event, error) {
	lggr := h.logger
	lggr.Infof("Handling message [%s]", event.Type)

	globalID, err := uuid.Parse(event.CorrelationID)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error parsing correlation ID")
		return nil, err
	}
	var req request
	err = parseInput(event.Data, &req)
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
		errEvent := events.NewEvent("order_creation_failed", "orders", map[string]interface{}{}).WithCorrelationID(event.CorrelationID)
		return errEvent, nil
	}

	lggr.Infof("Successfully created order [%s]", createRes)
	res := map[string]interface{}{"id": createRes.ID}
	successEvt := events.NewEvent("order_created", "orders", res).WithCorrelationID(event.CorrelationID)
	return successEvt, nil
}
