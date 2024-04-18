package workflows

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

type createOrderRequest struct {
	CustomerID   string `mapstructure:"customer_id"`
	Date         string `mapstructure:"date"`
	Total        *int64 `mapstructure:"total"`
	CurrencyCode string `mapstructure:"currency_code"`
}

type CreateOrderPayloadBuilder struct {
	logger *zap.SugaredLogger
}

func (b *CreateOrderPayloadBuilder) Build(ctx context.Context, data map[string]interface{}, action saga.ActionType) (map[string]interface{}, error) {
	lggr := b.logger
	// example of how to decode input data
	var req createOrderRequest
	err := mapstructure.Decode(data, &req)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error decoding input data")
		return nil, err
	}

	// example of how to build a payload
	payload := map[string]interface{}{}
	err = mapstructure.Decode(req, &payload)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error encoding payload")
		return nil, err
	}

	return payload, nil
}

func NewCreateOrderV1(logger *zap.SugaredLogger) *saga.Workflow {
	return &saga.Workflow{
		ID:           uuid.MustParse("2ef23373-9c01-4603-be2f-8e80552eb9a4"),
		Name:         "create_order",
		ReplyChannel: "saga.create-order.v1.response",
		Steps: saga.NewStepList(
			&saga.StepData{
				ID:             uuid.MustParse("4a4578ff-3602-4ad0-b262-6827c6ebc985"),
				Name:           "create_order",
				ServiceName:    "order",
				Compensable:    true,
				PayloadBuilder: &CreateOrderPayloadBuilder{},
			},
			&saga.StepData{
				ID:             uuid.MustParse("22d7c4bb-e751-4b47-a7a0-903ee5d3996e"),
				Name:           "verify_customer",
				ServiceName:    "customer",
				Compensable:    false,
				PayloadBuilder: &CreateOrderPayloadBuilder{},
			},
		),
	}
}
