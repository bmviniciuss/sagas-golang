package usecases

import (
	"context"
	"errors"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RejectOrderRequest struct {
	GlobalID uuid.UUID
}

type RejectOrderUseCasePort interface {
	Execute(ctx context.Context, request RejectOrderRequest) error
}

type RejectOrder struct {
	logger *zap.SugaredLogger
	repo   repositories.Orders
}

func NewRejectOrder(logger *zap.SugaredLogger, repo repositories.Orders) RejectOrder {
	return RejectOrder{logger: logger, repo: repo}
}

func (co RejectOrder) Execute(ctx context.Context, request RejectOrderRequest) error {
	lggr := co.logger
	lggr.Info("Rejecting order [%s]", request.GlobalID)
	order, err := co.repo.Find(ctx, request.GlobalID)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error finding order")
		return err
	}
	if order == nil {
		lggr.Errorf("Order [%s] not found", request.GlobalID)
		return errors.New("order not found")
	}

	order.Reject()
	err = co.repo.UpdateStatus(ctx, order)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error updating orders status")
		return err
	}

	lggr.Infof("Order [%s] rejected", request.GlobalID)
	return nil
}
