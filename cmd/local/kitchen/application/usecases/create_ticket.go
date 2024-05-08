package usecases

import (
	"context"
	"errors"

	"github.com/bmviniciuss/sagas-golang/cmd/local/kitchen/application/repositories"
	"github.com/bmviniciuss/sagas-golang/cmd/local/kitchen/domain/entities"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CreateTicketRequest struct {
	GlobalID     uuid.UUID
	CustomerID   uuid.UUID
	Amount       int64
	CurrencyCode string
	Items        []CreateTicketItem
}

type CreateTicketItem struct {
	ID        uuid.UUID
	Quantity  int32
	UnitPrice int64
}

type CreateTicketResponse struct {
	ID uuid.UUID
}

type CreateTicketUseCasePort interface {
	Execute(ctx context.Context, request CreateTicketRequest) (CreateTicketResponse, error)
}

type CreateTicket struct {
	logger *zap.SugaredLogger
	repo   repositories.Ticket
}

func NewCreateTicket(logger *zap.SugaredLogger, repo repositories.Ticket) CreateTicket {
	return CreateTicket{logger: logger, repo: repo}
}

func (co CreateTicket) Execute(ctx context.Context, request CreateTicketRequest) (CreateTicketResponse, error) {
	lggr := co.logger
	lggr.Info("Creating new ticket")

	items := make([]entities.Item, len(request.Items))
	for i, item := range request.Items {
		items[i] = entities.NewItem(item.ID, item.Quantity, item.UnitPrice)
	}
	ticket := entities.NewTicket(request.CustomerID, request.GlobalID, request.Amount, request.CurrencyCode, items)
	if ticket.Amount > 1000 {
		lggr.Info("Amount is greater than 1000")
		return CreateTicketResponse{}, errors.New("amount is greater than 1000")
	}

	err := co.repo.Insert(ctx, ticket)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error creating ticket")
		return CreateTicketResponse{}, err
	}
	return CreateTicketResponse{ID: ticket.ID}, nil
}
