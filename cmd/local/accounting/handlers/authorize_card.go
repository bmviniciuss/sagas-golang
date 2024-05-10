package handlers

import (
	"context"
	"encoding/json"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application"
	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"go.uber.org/zap"
)

type AuthorizeCardHandler struct {
	logger *zap.SugaredLogger
}

var (
	_ application.MessageHandler = (*AuthorizeCardHandler)(nil)
)

type request struct {
	Card   string `json:"card"`
	Amount *int64 `json:"amount"`
}

func NewAuthorizeCardHandler(logger *zap.SugaredLogger) *AuthorizeCardHandler {
	return &AuthorizeCardHandler{
		logger: logger,
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

func (h *AuthorizeCardHandler) Handle(ctx context.Context, msg *saga.Message) (*saga.Message, error) {
	lggr := h.logger
	lggr.Infof("Handling message [%s]", msg.EventType.String())

	globalID := msg.GlobalID
	var req request
	err := parseInput(msg.EventData, &req)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error reading input")
		return nil, err
	}
	lggr.Infof("Authoring card [%s] for amount [%d]", req.Card, *req.Amount)

	if req.Card == "0000000000000000" {
		lggr.Error("Card not authorized for this purchase")
		resEventType := saga.NewReplyEventType(msg.EventType, saga.FailureActionType)
		replyMessage := saga.NewParticipantMessage(globalID, nil, nil, resEventType, msg.Saga.ReplyChannel)
		return replyMessage, nil
	}

	lggr.Infof("Successfully authorized card [%s] for amount [%d]", req.Card, *req.Amount)
	res := map[string]interface{}{}
	resEventType := saga.NewReplyEventType(msg.EventType, saga.SuccessActionType)
	replyMessage := saga.NewParticipantMessage(globalID, res, nil, resEventType, msg.Saga.ReplyChannel)
	return replyMessage, nil
}
