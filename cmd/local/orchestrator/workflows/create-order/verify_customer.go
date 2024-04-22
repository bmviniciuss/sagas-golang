package createorder

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"go.uber.org/zap"
)

type VerifyCustomerPayloadBuilder struct {
	logger *zap.SugaredLogger
}

func NewVerifyCustomerPayloadBuilder(logger *zap.SugaredLogger) *VerifyCustomerPayloadBuilder {
	return &VerifyCustomerPayloadBuilder{logger: logger}
}

func (b *VerifyCustomerPayloadBuilder) Build(ctx context.Context, data map[string]interface{}, action saga.ActionType) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
