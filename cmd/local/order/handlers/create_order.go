package handlers

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application/usecases"
	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"go.uber.org/zap"
)

type CreateOrderHandler struct {
	logger             *zap.SugaredLogger
	createOrderUseCase usecases.CreateOrderUseCasePort
}

var (
	_ application.MessageHandler = (*CreateOrderHandler)(nil)
)

func NewCreateOrderHandler(logger *zap.SugaredLogger, createOrderUseCase usecases.CreateOrderUseCasePort) *CreateOrderHandler {
	return &CreateOrderHandler{
		logger:             logger,
		createOrderUseCase: createOrderUseCase,
	}
}

func (h *CreateOrderHandler) Handle(ctx context.Context, msg *saga.Message) (*saga.Message, error) {
	lggr := h.logger
	lggr.Infof("Handling message [%s]", msg.EventType.String())

	globalID := msg.GlobalID
	// TODO: receive the order data from the message
	orderID, err := h.createOrderUseCase.Execute(ctx, usecases.CreateOrderRequest{
		GlobalID: globalID,
	})
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error creating order")
		return nil, err
	}

	lggr.Infof("Successfully created order [%s]", orderID)
	replyMessage := saga.NewParticipantMessage(globalID, nil, nil, saga.SuccessActionType, msg)
	return replyMessage, nil
}
