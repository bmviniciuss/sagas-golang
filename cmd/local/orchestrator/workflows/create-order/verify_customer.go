package createorder

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/bmviniciuss/sagas-golang/pkg/events"
	"github.com/bmviniciuss/sagas-golang/pkg/structs"
	"go.uber.org/zap"
)

type VerifyCustomerPayloadBuilder struct {
	logger *zap.SugaredLogger
}

func NewVerifyCustomerPayloadBuilder(logger *zap.SugaredLogger) *VerifyCustomerPayloadBuilder {
	return &VerifyCustomerPayloadBuilder{logger: logger}
}

func (v *VerifyCustomerPayloadBuilder) Build(
	ctx context.Context,
	exec *saga.Execution,
	action saga.ActionType,
) (map[string]interface{}, error) {
	lggr := v.logger
	if action.IsRequest() {
		return v.buildRequestPayload(ctx, exec)
	}
	lggr.Infof("No payload to build for action: %s", action.String())
	return nil, nil
}

func (v *VerifyCustomerPayloadBuilder) buildRequestPayload(_ context.Context, exec *saga.Execution) (map[string]interface{}, error) {
	lggr := v.logger
	lggr.Info("Building request payload for verify customer step request")
	var input Input
	err := exec.Read("input", &input)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error decoding input data from execution in verify customer request builder")
		return nil, err
	}
	payload := VerifyCustomerRequestPayload{
		CustomerID: input.CustomerID,
	}
	lggr.Infof("Built request payload: %+v", payload)
	payloadMap, err := structs.ToMap(payload)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error converting payload to map")
		return nil, err
	}
	evt := events.NewEvent("verify_customer", "orchestrator", payloadMap).WithCorrelationID(exec.ID.String())
	evtMap, err := structs.ToMap(evt)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error converting event to map")
		return nil, err
	}
	return evtMap, nil
}
