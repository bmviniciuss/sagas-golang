package workflows

import (
	createorder "github.com/bmviniciuss/sagas-golang/cmd/local/orchestrator/workflows/create-order"
	"github.com/bmviniciuss/sagas-golang/internal/saga"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func NewCreateOrderV1(logger *zap.SugaredLogger) *saga.Workflow {
	return &saga.Workflow{
		ID:           uuid.MustParse("2ef23373-9c01-4603-be2f-8e80552eb9a4"), // TODO: remove uuids
		Name:         "create_order",
		ReplyChannel: "saga.create-order.v1.response",
		Steps: saga.NewStepList(
			&saga.StepData{
				ID:             uuid.MustParse("4a4578ff-3602-4ad0-b262-6827c6ebc985"),
				Name:           "create_order",
				ServiceName:    "order",
				Compensable:    true,
				PayloadBuilder: createorder.NewCreateOrderStepPayloadBuilder(logger),
			},
			&saga.StepData{
				ID:             uuid.MustParse("22d7c4bb-e751-4b47-a7a0-903ee5d3996e"),
				Name:           "verify_customer",
				ServiceName:    "customer",
				Compensable:    false,
				PayloadBuilder: createorder.NewVerifyCustomerPayloadBuilder(logger),
			},
		),
	}
}
