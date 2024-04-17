package repositories

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/domain/entities"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/presentation"
)

type Orders interface {
	List(ctx context.Context) ([]presentation.Order, error) // TODO: add pagination filters
	Insert(ctx context.Context, order entities.Order) error
}
