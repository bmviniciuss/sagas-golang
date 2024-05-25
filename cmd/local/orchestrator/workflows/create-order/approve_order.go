package createorder

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/bmviniciuss/sagas-golang/pkg/events"
	"github.com/bmviniciuss/sagas-golang/pkg/structs"
	"go.uber.org/zap"
)

type ApproveOrderPayloadBuilder struct {
	logger *zap.SugaredLogger
}

func NewApproveOrderPayloadBuilder(logger *zap.SugaredLogger) *ApproveOrderPayloadBuilder {
	return &ApproveOrderPayloadBuilder{logger: logger}
}

func (b *ApproveOrderPayloadBuilder) Build(ctx context.Context, exec *saga.Execution, action saga.ActionType) (map[string]interface{}, error) {
	lggr := b.logger
	if action.IsRequest() {
		return b.buildRequestPayload(ctx, exec)
	}
	lggr.Infof("No payload to build for action: %s", action.String())
	return nil, nil
}

func (b *ApproveOrderPayloadBuilder) buildRequestPayload(_ context.Context, exec *saga.Execution) (map[string]interface{}, error) {
	lggr := b.logger
	lggr.Info("Building request payload for approve order step")

	var input Input
	err := exec.Read("input", &input)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error decoding input data from execution in approve order step builder")
		return nil, err
	}

	evt := events.NewEvent("approve_order", "orchestrator", map[string]interface{}{}).WithCorrelationID(exec.ID.String())
	eventMap, err := structs.ToMap(evt)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error converting event to map")
		return nil, err
	}

	return eventMap, nil
}
