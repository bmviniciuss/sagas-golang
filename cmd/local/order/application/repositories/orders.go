package repositories

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/domain/entities"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/presentation"
	"github.com/google/uuid"
)

type Orders interface {
	List(ctx context.Context) ([]presentation.Order, error) // TODO: add pagination filters
	Insert(ctx context.Context, order entities.Order) error
	FindByID(ctx context.Context, id uuid.UUID) (*presentation.OrderById, error)
}
