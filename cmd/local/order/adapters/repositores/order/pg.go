package order

import (
	"context"
	"time"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application/repositories"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/domain/entities"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/presentation"
	"github.com/bmviniciuss/sagas-golang/pkg/utc"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
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
	_ repositories.Orders = (*RepositoryAdapter)(nil)
)

type orderModel struct {
	ID           string    `db:"id"`
	UUID         string    `db:"uuid"`
	GlobalID     string    `db:"global_id"`
	ClientID     string    `db:"client_id"`
	CustomerID   string    `db:"customer_id"`
	Total        int64     `db:"total"`
	CurrencyCode string    `db:"currency_code"`
	Status       string    `db:"status"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

const listOrdersQuery = `
SELECT id, uuid, global_id, client_id, customer_id, total, currency_code, status, created_at, updated_at
FROM orders.orders
`

func (r *RepositoryAdapter) List(ctx context.Context) ([]presentation.Order, error) {
	lggr := r.lggr
	lggr.Info("RepositoryAdapter.List")
	db, err := r.pool.Acquire(ctx)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error acquiring connection")
		return nil, err
	}
	defer db.Release()
	rows, err := db.Query(ctx, listOrdersQuery)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []presentation.Order{}, nil
		}
		lggr.With(zap.Error(err)).Error("Got error querying database")
		return nil, err
	}
	var orders []orderModel
	err = pgxscan.ScanAll(&orders, rows)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error scanning rows")
		return nil, err
	}
	ordersPresentation := make([]presentation.Order, len(orders))
	for i, order := range orders {
		ordersPresentation[i] = presentation.Order{
			ID:           order.UUID,
			GlobalID:     order.GlobalID,
			ClientID:     order.ClientID,
			CustomerID:   order.CustomerID,
			Total:        order.Total,
			CurrencyCode: order.CurrencyCode,
			Status:       order.Status,
			CreatedAt:    utc.NewFromTime(order.CreatedAt),
			UpdatedAt:    utc.NewFromTime(order.UpdatedAt),
		}
	}
	return ordersPresentation, nil
}

const insertOrderQuery = `
INSERT INTO orders.orders (uuid, global_id, client_id, customer_id, total, currency_code, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`

func (r *RepositoryAdapter) Insert(ctx context.Context, order entities.Order) error {
	lggr := r.lggr
	lggr.Info("RepositoryAdapter.Insert")

	db, err := r.pool.Acquire(ctx)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error acquiring connection")
		return err
	}
	defer db.Release()

	_, err = db.Exec(ctx, insertOrderQuery,
		order.ID.String(),
		order.GlobalID.String(),
		order.ClientID.String(),
		order.CustomerID.String(),
		order.Total,
		order.CurrencyCode,
		order.Status.String(),
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error inserting order")
		return err
	}

	return nil
}
