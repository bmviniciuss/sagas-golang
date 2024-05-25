package handlers

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application/usecases"
	"github.com/bmviniciuss/sagas-golang/pkg/events"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ApproveOrder struct {
	logger              *zap.SugaredLogger
	approveOrderUseCase usecases.ApproveOrderUseCasePort
}

var (
	_ application.MessageHandler = (*ApproveOrder)(nil)
)

func NewApproveOrder(logger *zap.SugaredLogger, approveOrderUseCase usecases.ApproveOrderUseCasePort) *ApproveOrder {
	return &ApproveOrder{
		logger:              logger,
		approveOrderUseCase: approveOrderUseCase,
	}
}

// TODO: add enum
// Request:      "approve_order",
// Success:      "order_approved",

func (h *ApproveOrder) Handle(ctx context.Context, event *events.Event) (*events.Event, error) {
	lggr := h.logger
	lggr.Infof("Handling message [%s]", event.Type)

	globalID, err := uuid.Parse(event.CorrelationID)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error parsing correlation ID")
		return nil, err
	}

	err = h.approveOrderUseCase.Execute(ctx, usecases.ApproveOrderRequest{
		GlobalID: globalID,
	})
	if err != nil {
		lggr.Info("Got error approving order. Because this is a retryable operation, after a go-non-go decision, the event should be sent to retry my any mechanism")
		return nil, nil
	}
	successEvt := events.NewEvent("order_approved", "orders", map[string]interface{}{}).WithCorrelationID(event.CorrelationID)
	return successEvt, nil
}
