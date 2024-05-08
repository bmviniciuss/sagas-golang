package handlers

import (
	"context"
	"encoding/json"

	"github.com/bmviniciuss/sagas-golang/cmd/local/kitchen/application"
	"github.com/bmviniciuss/sagas-golang/cmd/local/kitchen/application/usecases"
	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CreateTicketHandler struct {
	logger              *zap.SugaredLogger
	createTicketUseCase usecases.CreateTicketUseCasePort
}

var (
	_ application.MessageHandler = (*CreateTicketHandler)(nil)
)

func NewCreateTicketHandler(logger *zap.SugaredLogger, createTicketUseCase usecases.CreateTicketUseCasePort) *CreateTicketHandler {
	return &CreateTicketHandler{
		logger:              logger,
		createTicketUseCase: createTicketUseCase,
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

type request struct {
	CustomerID   uuid.UUID `json:"customer_id"`
	Amount       *int64    `json:"amount"`
	CurrencyCode string    `json:"currency_code"`
	Items        []item    `json:"items"`
}

type item struct {
	ID        uuid.UUID `json:"id"`
	Quantity  *int32    `json:"quantity"`
	UnitPrice *int64    `json:"unit_price"`
}

func (h *CreateTicketHandler) Handle(ctx context.Context, msg *saga.Message) (*saga.Message, error) {
	lggr := h.logger
	lggr.Infof("Handling message [%s]", msg.EventType.String())

	globalID := msg.GlobalID
	var req request
	err := parseInput(msg.EventData, &req)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error reading input")
		return nil, err
	}

	items := make([]usecases.CreateTicketItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = usecases.CreateTicketItem{
			ID:        item.ID,
			Quantity:  *item.Quantity,
			UnitPrice: *item.UnitPrice,
		}
	}

	createTicketResponse, err := h.createTicketUseCase.Execute(ctx, usecases.CreateTicketRequest{
		GlobalID:     globalID,
		CustomerID:   req.CustomerID,
		Amount:       *req.Amount,
		CurrencyCode: req.CurrencyCode,
		Items:        items,
	})
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error creating ticket")
		resEventType := saga.NewReplyEventType(msg.EventType, saga.FailureActionType)
		replyMessage := saga.NewParticipantMessage(globalID, nil, nil, resEventType, msg.Saga.ReplyChannel)
		return replyMessage, nil
	}

	lggr.Infof("Successfully created ticket [%s]", createTicketResponse)
	res := map[string]interface{}{"ticket_id": createTicketResponse.ID}
	resEventType := saga.NewReplyEventType(msg.EventType, saga.SuccessActionType)
	replyMessage := saga.NewParticipantMessage(globalID, res, nil, resEventType, msg.Saga.ReplyChannel)
	return replyMessage, nil
}
