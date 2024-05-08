package ticket

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/cmd/local/kitchen/adapters/repositores/ticket/generated"
	"github.com/bmviniciuss/sagas-golang/cmd/local/kitchen/application/repositories"
	"github.com/bmviniciuss/sagas-golang/cmd/local/kitchen/domain/entities"
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
