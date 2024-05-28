package handlers

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application/usecases"
	"github.com/bmviniciuss/sagas-golang/pkg/events"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RejectOrder struct {
	logger             *zap.SugaredLogger
	rejectOrderUseCase usecases.RejectOrderUseCasePort
}

var (
	_ application.MessageHandler = (*ApproveOrder)(nil)
)

func NewRejectOrder(logger *zap.SugaredLogger, rejectOrderUseCase usecases.RejectOrderUseCasePort) *RejectOrder {
	return &RejectOrder{
		logger:             logger,
		rejectOrderUseCase: rejectOrderUseCase,
	}
}

// TODO: add enum
// reject_order
// order_rejected

func (h *RejectOrder) Handle(ctx context.Context, event *events.Event) (*events.Event, error) {
	lggr := h.logger
	lggr.Infof("Handling message [%s]", event.Type)

	globalID, err := uuid.Parse(event.CorrelationID)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error parsing correlation ID")
		return nil, err
	}

	err = h.rejectOrderUseCase.Execute(ctx, usecases.RejectOrderRequest{
		GlobalID: globalID,
	})
	if err != nil {
		lggr.Info("Got error rejecting order. Because this is a compesation operation, after a go-non-go decision, the event should be sent to retry my any mechanism")
		return nil, nil
	}
	successEvt := events.NewEvent("order_rejected", "orders", map[string]interface{}{}).WithCorrelationID(event.CorrelationID)
	return successEvt, nil
}
