package ticket

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/cmd/local/kitchen/adapters/repositores/ticket/generated"
	"github.com/bmviniciuss/sagas-golang/cmd/local/kitchen/application/repositories"
	"github.com/bmviniciuss/sagas-golang/cmd/local/kitchen/domain/entities"
	"github.com/bmviniciuss/sagas-golang/pkg/utc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type RepositoryAdapter struct {
	lggr *zap.SugaredLogger
	pool *pgxpool.Pool
}

func NewRepositoryAdapter(lggr *zap.SugaredLogger, pool *pgxpool.Pool) *RepositoryAdapter {
	return &RepositoryAdapter{lggr: lggr, pool: pool}
}

var (
	_ repositories.Ticket = (*RepositoryAdapter)(nil)
)

func (r *RepositoryAdapter) Insert(ctx context.Context, ticket *entities.Ticket) error {
	lggr := r.lggr
	lggr.Info("RepositoryAdapter.Insert")

	db, err := r.pool.Acquire(ctx)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error acquiring connection")
		return err
	}
	defer db.Release()
	queries := generated.New(db)

	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error starting transaction")
		return err
	}
	defer tx.Rollback(ctx)
	qtx := queries.WithTx(tx)

	id, err := qtx.InsertTicket(ctx, generated.InsertTicketParams{
		Uuid:         ticket.ID,
		CustomerID:   ticket.CustomerID,
		Status:       ticket.Status.String(),
		Amount:       ticket.Amount,
		CurrencyCode: ticket.CurrencyCode,
		CreatedAt: pgtype.Timestamptz{
			Time:  ticket.CreatedAt.Time(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  ticket.UpdatedAt.Time(),
			Valid: true,
		},
	})
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error inserting ticket")
		return err
	}

	for _, item := range ticket.Items {
		err = qtx.InsertTicketItem(ctx, generated.InsertTicketItemParams{
			Uuid:      item.ID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			TicketID:  id,
		})
		if err != nil {
			lggr.With(zap.Error(err)).Error("Got error inserting ticket item")
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error committing transaction")
		return err
	}

	return nil
}

func (r *RepositoryAdapter) Find(ctx context.Context, ticketID uuid.UUID) (*entities.Ticket, error) {
	lggr := r.lggr
	lggr.Info("RepositoryAdapter.Find")

	db, err := r.pool.Acquire(ctx)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error acquiring connection")
		return nil, err
	}
	defer db.Release()
	queries := generated.New(db)

	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error starting transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := queries.WithTx(tx)

	ticketRow, err := qtx.GetTicket(ctx, ticketID)
	if err != nil {
		if err == pgx.ErrNoRows {
			lggr.Error("Ticket not found")
			return &entities.Ticket{}, nil
		}
		lggr.With(zap.Error(err)).Error("Got error")
		return nil, err
	}

	itemsRows, err := qtx.GetTicketItems(ctx, ticketRow.Identifier)
	if err != nil && err != pgx.ErrNoRows {
		lggr.With(zap.Error(err)).Error("Got error getting ticket items")
		return nil, err
	}

	items := make([]entities.Item, len(itemsRows))
	for i, itemRow := range itemsRows {
		items[i] = entities.Item{
			ID:        itemRow.Uuid,
			Quantity:  itemRow.Quantity,
			UnitPrice: itemRow.UnitPrice,
		}
	}

	ticket := &entities.Ticket{
		ID:           ticketRow.Uuid,
		CustomerID:   ticketRow.CustomerID,
		Amount:       ticketRow.Amount,
		CurrencyCode: ticketRow.CurrencyCode,
		Status:       entities.TicketStatus(ticketRow.Status),
		CreatedAt:    utc.NewFromTime(ticketRow.CreatedAt.Time),
		UpdatedAt:    utc.NewFromTime(ticketRow.UpdatedAt.Time),
		Items:        items,
	}

	err = tx.Commit(ctx)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error committing transaction")
		return nil, err
	}

	return ticket, nil
}

func (r *RepositoryAdapter) UpdateStatus(ctx context.Context, ticket *entities.Ticket) error {
	lggr := r.lggr
	lggr.Info("RepositoryAdapter.UpdateStatus")

	db, err := r.pool.Acquire(ctx)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error acquiring connection")
		return err
	}
	defer db.Release()
	queries := generated.New(db)

	err = queries.UpdateTicketStatus(ctx, generated.UpdateTicketStatusParams{
		Uuid:   ticket.ID,
		Status: ticket.Status.String(),
		UpdatedAt: pgtype.Timestamptz{
			Time:  ticket.UpdatedAt.Time(),
			Valid: true,
		},
	})
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error updating ticket status")
		return err
	}

	lggr.Infof("Ticket [%s] status updated to [%s]", ticket.ID.String(), ticket.Status.String())
	return nil
}
