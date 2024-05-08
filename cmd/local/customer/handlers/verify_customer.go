package handlers

import (
	"context"
	"encoding/json"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application"
	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type VerifyCustomer struct {
	logger *zap.SugaredLogger
}

var (
	_ application.MessageHandler = (*VerifyCustomer)(nil)
)

type request struct {
	CustomerID uuid.UUID `json:"customer_id"`
}

func NewVerifyCustomer(logger *zap.SugaredLogger) *VerifyCustomer {
	return &VerifyCustomer{
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

func (h *VerifyCustomer) Handle(ctx context.Context, msg *saga.Message) (*saga.Message, error) {
	lggr := h.logger
	lggr.Infof("Handling message [%s]", msg.EventType.String())

	globalID := msg.GlobalID
	var req request
	err := parseInput(msg.EventData, &req)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error reading input")
		return nil, err
	}
	lggr.Infof("Validating customer [%s]", req.CustomerID.String())

	if req.CustomerID.String() == "00000000-0000-0000-0000-000000000000" {
		lggr.Error("Customer not available to create order")
		resEventType := saga.NewReplyEventType(msg.EventType, saga.FailureActionType)
		replyMessage := saga.NewParticipantMessage(globalID, nil, nil, resEventType, msg.Saga.ReplyChannel)
		return replyMessage, nil
	}
	// if err != nil {
	// 	lggr.With(zap.Error(err)).Error("Got error creating order")
	// 	resEventType := saga.NewReplyEventType(msg.EventType, saga.FailureActionType)
	// 	replyMessage := saga.NewParticipantMessage(globalID, nil, nil, resEventType, msg.Saga.ReplyChannel)
	// 	return replyMessage, nil
	// }

	lggr.Infof("Customer can create order")
	res := map[string]interface{}{"customer_id": req.CustomerID.String()}
	resEventType := saga.NewReplyEventType(msg.EventType, saga.SuccessActionType)
	replyMessage := saga.NewParticipantMessage(globalID, res, nil, resEventType, msg.Saga.ReplyChannel)
	return replyMessage, nil
}
