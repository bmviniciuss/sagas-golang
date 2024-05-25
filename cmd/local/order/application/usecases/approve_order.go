package usecases

import (
	"context"
	"errors"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ApproveOrderRequest struct {
	GlobalID uuid.UUID
}

type ApproveOrderUseCasePort interface {
	Execute(ctx context.Context, request ApproveOrderRequest) error
}

type ApproveOrder struct {
	logger *zap.SugaredLogger
	repo   repositories.Orders
}

func NewApproveOrder(logger *zap.SugaredLogger, repo repositories.Orders) ApproveOrder {
	return ApproveOrder{logger: logger, repo: repo}
}

func (co ApproveOrder) Execute(ctx context.Context, request ApproveOrderRequest) error {
	lggr := co.logger
	lggr.Info("Approving order [%s]", request.GlobalID)
	order, err := co.repo.Find(ctx, request.GlobalID)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error finding order")
		return err
	}
	if order == nil {
		lggr.Errorf("Order [%s] not found", request.GlobalID)
		return errors.New("order not found")
	}

	order.Approve()
	err = co.repo.UpdateStatus(ctx, order)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error updating orders status")
		return err
	}

	lggr.Infof("Order [%s] approved", request.GlobalID)
	return nil
}
