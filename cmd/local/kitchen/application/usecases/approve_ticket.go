package usecases

import (
	"context"
	"errors"

	"github.com/bmviniciuss/sagas-golang/cmd/local/kitchen/application/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ApproveTicketRequest struct {
	TicketID uuid.UUID
}

type ApproveTicketUseCasePort interface {
	Execute(ctx context.Context, request ApproveTicketRequest) error
}

type ApproveTicket struct {
	logger *zap.SugaredLogger
	repo   repositories.Ticket
}

func NewApproveTicket(logger *zap.SugaredLogger, repo repositories.Ticket) ApproveTicket {
	return ApproveTicket{logger: logger, repo: repo}
}

func (co ApproveTicket) Execute(ctx context.Context, request ApproveTicketRequest) error {
	lggr := co.logger
	lggr.Infof("Accepting ticket [%s]", request.TicketID.String())

	ticket, err := co.repo.Find(ctx, request.TicketID)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error finding ticket")
		return err
	}
	if ticket.IsEmpty() {
		lggr.Error("Ticket not found")
		return errors.New("ticket not found")
	}

	ticket.Approve()
	err = co.repo.UpdateStatus(ctx, ticket)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error updating ticket status")
		return err
	}
	return nil
}
