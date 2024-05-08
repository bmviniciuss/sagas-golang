package usecases

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application/repositories"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/presentation"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type GetOrderByID struct {
	lggr             *zap.SugaredLogger
	ordersRepository repositories.Orders
}

func NewGetOrderByID(lggr *zap.SugaredLogger, ordersRepository repositories.Orders) *GetOrderByID {
	return &GetOrderByID{
		lggr:             lggr,
		ordersRepository: ordersRepository,
	}
}

func (l *GetOrderByID) Execute(ctx context.Context, id uuid.UUID) (*presentation.OrderById, error) {
	l.lggr.Info("GetOrderByID.Execute")
	order, err := l.ordersRepository.FindByID(ctx, id)
	if err != nil {
		l.lggr.With(zap.Error(err)).Error("Got error listing orders")
		return nil, err
	}
	return order, nil
}
