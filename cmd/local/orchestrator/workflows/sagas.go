package workflows

import (
	createorder "github.com/bmviniciuss/sagas-golang/cmd/local/orchestrator/workflows/create-order"
	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"go.uber.org/zap"
)

func NewCreateOrderV1(logger *zap.SugaredLogger) *saga.Workflow {
	return &saga.Workflow{
		Name:         "create_order_v1",
		ReplyChannel: "saga.create_order_v1.response",
		Steps: saga.NewStepList(
			&saga.StepData{
				Name:           "create_order",
				ServiceName:    "order",
				Compensable:    true,
				PayloadBuilder: createorder.NewCreateOrderStepPayloadBuilder(logger),
			},
			&saga.StepData{
				Name:           "verify_customer",
				ServiceName:    "customer",
				Compensable:    false,
				PayloadBuilder: createorder.NewVerifyCustomerPayloadBuilder(logger),
			},
		),
	}
}
