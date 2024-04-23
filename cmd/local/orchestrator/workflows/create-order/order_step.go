package createorder

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/zap"
)

type CreateOrderStepPayloadBuilder struct {
	logger *zap.SugaredLogger
}

func NewCreateOrderStepPayloadBuilder(logger *zap.SugaredLogger) *CreateOrderStepPayloadBuilder {
	return &CreateOrderStepPayloadBuilder{logger: logger}
}

func (b *CreateOrderStepPayloadBuilder) Build(ctx context.Context, data map[string]interface{}, action saga.ActionType) (map[string]interface{}, error) {
	// TODO: add execution as parameter
	lggr := b.logger
	// example of how to decode input data
	var req Input
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
