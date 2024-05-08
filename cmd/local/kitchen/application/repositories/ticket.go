package repositories

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/cmd/local/kitchen/domain/entities"
)

type Ticket interface {
	Insert(ctx context.Context, ticket *entities.Ticket) error
}
