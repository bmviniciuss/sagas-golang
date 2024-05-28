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
				ServiceName:    "orders",
				Compensable:    true,
				PayloadBuilder: createorder.NewCreateOrderStepPayloadBuilder(logger),
				EventTypes: saga.EventTypes{
					Request:            "create_order",
					CompesationRequest: "reject_order",
					Success:            "order_created",
					Failure:            "order_creation_failed",
					Compensation:       "order_rejected",
				},
				Topics: saga.Topics{
					Request:  "service.orders.request",
					Response: "service.orders.events",
				},
			},
			&saga.StepData{
				Name:           "verify_customer",
				ServiceName:    "customers",
				Compensable:    false,
				PayloadBuilder: createorder.NewVerifyCustomerPayloadBuilder(logger),
				EventTypes: saga.EventTypes{
					Request:      "verify_customer",
					Success:      "customer_verified",
					Failure:      "customer_verification_failed",
					Compensation: "",
				},
				Topics: saga.Topics{
					Request:  "service.customers.request",
					Response: "service.customers.events",
				},
			},
			&saga.StepData{
				Name:           "authorize_card",
				ServiceName:    "accounting",
				Compensable:    false,
				PayloadBuilder: createorder.NewAuthorizeCardPayloadBuilder(logger),
				EventTypes: saga.EventTypes{
					Request:      "authorize_card",
					Success:      "card_authorized",
					Failure:      "card_authorization_failed",
					Compensation: "",
				},
				Topics: saga.Topics{
					Request:  "service.accounting.request",
					Response: "service.accounting.events",
				},
			},
			&saga.StepData{
				Name:           "approve_order",
				ServiceName:    "orders",
				Compensable:    true,
				PayloadBuilder: createorder.NewApproveOrderPayloadBuilder(logger),
				EventTypes: saga.EventTypes{
					Request:      "approve_order",
					Success:      "order_approved",
					Failure:      "",
					Compensation: "",
				},
				Topics: saga.Topics{
					Request:  "service.orders.request",
					Response: "service.orders.events",
				},
			},
		),
	}
}
