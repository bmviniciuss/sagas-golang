package handlers

import (
	"context"
	"encoding/json"

	"github.com/bmviniciuss/sagas-golang/cmd/local/accounting/application"
	"github.com/bmviniciuss/sagas-golang/pkg/events"
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

// TODO: add enum
// Request:      "authorize_card",
// Success:      "card_authorized",
// Failure:      "card_authorization_failed",

func (h *AuthorizeCardHandler) Handle(ctx context.Context, msg *events.Event) (*events.Event, error) {
	lggr := h.logger
	lggr.Infof("Handling message [%s]", msg.Type)

	globalID := msg.CorrelationID
	var req request
	err := parseInput(msg.Data, &req)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error reading input")
		return nil, err
	}
	lggr.Infof("Authoring card [%s] for amount [%d]", req.Card, *req.Amount)

	if req.Card == "0000000000000000" {
		lggr.Error("Card not authorized for this purchase")
		errEvt := events.NewEvent("card_authorization_failed", "accounting", make(map[string]interface{})).WithCorrelationID(globalID)
		return errEvt, nil
	}

	lggr.Infof("Successfully authorized card [%s] for amount [%d]", req.Card, *req.Amount)
	res := map[string]interface{}{}
	successEvt := events.NewEvent("card_authorized", "accounting", res).WithCorrelationID(globalID)
	return successEvt, nil
}
