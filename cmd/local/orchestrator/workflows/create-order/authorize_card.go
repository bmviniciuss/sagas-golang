package createorder

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/bmviniciuss/sagas-golang/pkg/structs"
	"go.uber.org/zap"
)

type AuthorizeCardPayloadBuilder struct {
	logger *zap.SugaredLogger
}

func NewAuthorizeCardPayloadBuilder(logger *zap.SugaredLogger) *AuthorizeCardPayloadBuilder {
	return &AuthorizeCardPayloadBuilder{logger: logger}
}

func (v *AuthorizeCardPayloadBuilder) Build(
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

func (v *AuthorizeCardPayloadBuilder) buildRequestPayload(_ context.Context, exec *saga.Execution) (map[string]interface{}, error) {
	lggr := v.logger
	lggr.Info("Building request payload for card authorization step request")
	var input Input
	err := exec.Read("input", &input)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error reading input")
		return nil, err
	}
	payload := AuthorizeCardRequestPayload{
		Card:   input.Card,
		Amount: *input.Amount,
	}
	lggr.Infof("Built request payload: %+v", payload)
	payloadMap, err := structs.ToMap(payload)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error converting payload to map")
		return nil, err
	}
	return payloadMap, nil
}