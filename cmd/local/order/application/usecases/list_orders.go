package usecases

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application/repositories"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/presentation"
	"go.uber.org/zap"
)

type ListOrders struct {
	lggr             *zap.SugaredLogger
	ordersRepository repositories.Orders
}

func NewListOrders(lggr *zap.SugaredLogger, ordersRepository repositories.Orders) *ListOrders {
	return &ListOrders{
		lggr:             lggr,
		ordersRepository: ordersRepository,
	}
}

type ListOrderResults struct {
	Orders []presentation.Order `json:"orders"`
}

func (l *ListOrders) Execute(ctx context.Context) (ListOrderResults, error) {
	l.lggr.Info("ListOrders.Execute")
	orders, err := l.ordersRepository.List(ctx)
	if err != nil {
		l.lggr.With(zap.Error(err)).Error("Got error listing orders")
		return ListOrderResults{}, err
	}
	return ListOrderResults{Orders: orders}, nil
}
