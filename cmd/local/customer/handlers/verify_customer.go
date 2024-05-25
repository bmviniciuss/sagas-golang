package handlers

import (
	"context"
	"encoding/json"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application"
	"github.com/bmviniciuss/sagas-golang/pkg/events"
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

// Request:      "verify_customer",
// Success:      "customer_verified",
// Failure:      "customer_verification_failed",

func (h *VerifyCustomer) Handle(ctx context.Context, msg *events.Event) (*events.Event, error) {
	lggr := h.logger
	lggr.Infof("Handling message with message [%s]", msg.Type)

	globalID := msg.CorrelationID
	var req request
	err := parseInput(msg.Data, &req)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error reading input")
		return nil, err
	}
	lggr.Infof("Validating customer [%s]", req.CustomerID.String())

	if req.CustomerID.String() == "00000000-0000-0000-0000-000000000000" {
		lggr.Error("Customer not available to create order")
		errEvt := events.NewEvent("customer_verification_failed", "customers", nil).WithCorrelationID(globalID)
		return errEvt, nil
	}

	lggr.Infof("Customer can create order")
	res := map[string]interface{}{"customer_id": req.CustomerID.String()}
	successEvt := events.NewEvent("customer_verified", "customers", res).WithCorrelationID(globalID)
	return successEvt, nil
}
