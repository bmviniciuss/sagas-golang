package repositories

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/cmd/local/kitchen/domain/entities"
	"github.com/google/uuid"
)

type Ticket interface {
	Insert(ctx context.Context, ticket *entities.Ticket) error
	Find(ctx context.Context, ticketID uuid.UUID) (*entities.Ticket, error)
	UpdateStatus(ctx context.Context, ticket *entities.Ticket) error
}
