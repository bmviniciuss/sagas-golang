package order

import (
	"context"

	"github.com/bmviniciuss/sagas-golang/cmd/local/order/adapters/repositores/order/generated"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/application/repositories"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/domain/entities"
	"github.com/bmviniciuss/sagas-golang/cmd/local/order/presentation"
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
	_ repositories.Orders = (*RepositoryAdapter)(nil)
)

func (r *RepositoryAdapter) List(ctx context.Context) ([]presentation.Order, error) {
	lggr := r.lggr
	lggr.Info("RepositoryAdapter.List")
	db, err := r.pool.Acquire(ctx)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error acquiring connection")
		return nil, err
	}
	defer db.Release()
	queries := generated.New(db)
	rows, err := queries.ListOrders(ctx)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []presentation.Order{}, nil
		}
		lggr.With(zap.Error(err)).Error("Got error querying database")
		return nil, err
	}

	ordersPresentation := make([]presentation.Order, len(rows))
	for i, row := range rows {
		ordersPresentation[i] = presentation.Order{
			ID:           row.Uuid.String(),
			CustomerID:   row.CustomerID.String(),
			Amount:       row.Amount,
			CurrencyCode: row.CurrencyCode,
			Status:       row.Status,
			CreatedAt:    utc.NewFromTime(row.CreatedAt.Time),
			UpdatedAt:    utc.NewFromTime(row.UpdatedAt.Time),
		}
	}
	return ordersPresentation, nil
}

func (r *RepositoryAdapter) Insert(ctx context.Context, order entities.Order) error {
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
		lggr.With(zap.Error(err)).Error("Got error beginning transaction")
		return err
	}
	defer tx.Rollback(ctx)
	qtx := queries.WithTx(tx)

	orderID, err := qtx.InsertOrder(ctx, generated.InsertOrderParams{
		Uuid:         order.ID,
		CustomerID:   order.CustomerID,
		Status:       order.Status.String(),
		Amount:       order.Amount,
		CurrencyCode: order.CurrencyCode,
		CreatedAt: pgtype.Timestamptz{
			Time:  order.CreatedAt.Time(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time:  order.UpdatedAt.Time(),
			Valid: true,
		},
	})
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error inserting order")
		return err
	}

	for _, itemRow := range order.Items {
		err = qtx.InsertOrderItem(ctx, generated.InsertOrderItemParams{
			Uuid:      itemRow.ID,
			Quantity:  itemRow.Quantity,
			UnitPrice: itemRow.UnitPrice,
			OrderID:   orderID,
		})
		if err != nil {
			lggr.With(zap.Error(err)).Error("Got error inserting order item")
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error committing insert order transaction")
		return err
	}

	return nil
}

func (r *RepositoryAdapter) FindByID(ctx context.Context, id uuid.UUID) (*presentation.OrderById, error) {
	lggr := r.lggr
	lggr.Info("RepositoryAdapter.Insert")

	db, err := r.pool.Acquire(ctx)
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error acquiring connection")
		return nil, err
	}
	defer db.Release()

	queries := generated.New(db)

	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		lggr.With(zap.Error(err)).Error("Got error beginning transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := queries.WithTx(tx)

	order, err := qtx.GetOrder(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &presentation.OrderById{}, nil
		}
		lggr.With(zap.Error(err)).Error("Got error querying get order by id")
		return nil, err
	}

	orderItems := []presentation.Item{}
	orderItemsRows, err := qtx.GetOrderItems(ctx, order.Uuid)
	if err != nil && err != pgx.ErrNoRows {
		lggr.With(zap.Error(err)).Error("Got error querying get order items by order id")
		return nil, err
	}

	for _, item := range orderItemsRows {
		orderItems = append(orderItems, presentation.Item{
			ID:        item.Uuid.String(),
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})
	}

	orderPresentation := presentation.OrderById{
		ID:           order.Uuid.String(),
		CustomerID:   order.CustomerID.String(),
		Amount:       order.Amount,
		CurrencyCode: order.CurrencyCode,
		Status:       order.Status,
		Items:        orderItems,
		CreatedAt:    utc.NewFromTime(order.CreatedAt.Time),
		UpdatedAt:    utc.NewFromTime(order.UpdatedAt.Time),
	}

	return &orderPresentation, nil
}
