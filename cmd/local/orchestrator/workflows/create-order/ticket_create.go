package createorder

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/bmviniciuss/sagas-golang/pkg/structs"
	"go.uber.org/zap"
)

type CreateTicketPayloadBuilder struct {
	logger *zap.SugaredLogger
}

func NewCreateTicketPayloadBuilder(logger *zap.SugaredLogger) *CreateTicketPayloadBuilder {
	return &CreateTicketPayloadBuilder{logger: logger}
}

func (v *CreateTicketPayloadBuilder) Build(
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

func (v *CreateTicketPayloadBuilder) buildRequestPayload(_ context.Context, exec *saga.Execution) (map[string]interface{}, error) {
	lggr := v.logger
	lggr.Info("Building request payload for create ticket message")
	var input Input
	err := exec.Read("input", &input)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error decoding input data from execution in verify customer request builder")
		return nil, err
	}

	items := make([]CreateTicketItemRequestPayload, len(input.Items))
	for i, item := range input.Items {
		items[i] = CreateTicketItemRequestPayload{
			ID:        item.ID,
			Quantity:  *item.Quantity,
			UnitPrice: *item.UnitPrice,
		}
	}
	payload := CreateTicketRequestPayload{
		CustomerID:   input.CustomerID,
		Amount:       *input.Amount,
		CurrencyCode: input.CurrencyCode,
		Items:        items,
	}
	lggr.Infof("Built request payload: %+v", payload)
	payloadMap, err := structs.ToMap(payload)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error converting payload to map")
		return nil, err
	}
	return payloadMap, nil
}
