package orderrepository

import (
	"context"
	"fmt"
	"route256/loms/internal/pkg/model"
	"route256/loms/internal/pkg/repository"
	"route256/loms/internal/pkg/repository/order_repository/sqlc_order"
	"route256/loms/pkg/tracing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

type DbOrderRepository struct {
	db      DB
	queries *sqlc_order.Queries
}

func NewDbOrderRepository(db DB) *DbOrderRepository {
	return &DbOrderRepository{
		db:      db,
		queries: sqlc_order.New(db),
	}
}

func (r *DbOrderRepository) Create(ctx context.Context, order model.Order, inTx ...func(context.Context, model.OrderID, model.OrderStatus) error) (_ model.OrderID, err error) {
	ctx, span := tracing.Start(ctx, "DbOrderRepository.Create")
	defer tracing.EndWithCheckError(span, &err)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("r.db.Begin: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	if order.Status == model.OrderStatusNone {
		order.Status = model.OrderStatusNew
	}

	id, err := qtx.Create(ctx, sqlc_order.CreateParams{
		UserID: int64(order.User),
		Status: string(order.Status),
	})
	if err != nil {
		return 0, fmt.Errorf("qtx.Create: %w", err)
	}

	for _, item := range order.Items {
		err := qtx.AddItem(ctx, sqlc_order.AddItemParams{
			OrderID: int64(id),
			Sku:     int64(item.Sku),
			Count:   int32(item.Count),
		})
		if err != nil {
			return 0, fmt.Errorf("qtx.AddItem: %w", err)
		}
	}

	ctxTx := context.WithValue(ctx, repository.CtxTxKey{}, tx)
	for _, f := range inTx {
		if err := f(ctxTx, model.OrderID(id), order.Status); err != nil {
			return 0, fmt.Errorf("inTx: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("tx.Commit: %w", err)
	}
	return model.OrderID(id), nil
}

func (r *DbOrderRepository) GetById(ctx context.Context, orderID model.OrderID) (_ model.Order, err error) {
	ctx, span := tracing.Start(ctx, "DbOrderRepository.GetById")
	defer tracing.EndWithCheckError(span, &err)

	orderItems, err := r.queries.GetById(ctx, int64(orderID))
	if err != nil {
		return model.Order{}, fmt.Errorf("r.queries.GetById: %w", err)
	}
	if len(orderItems) == 0 {
		return model.Order{}, fmt.Errorf("invalid order id: %d", orderID)
	}
	order := model.Order{
		Status: model.OrderStatus(orderItems[0].Order.Status),
		User:   model.UserID(orderItems[0].Order.UserID),
		Items:  make([]model.OrderItem, 0, len(orderItems)),
	}
	for _, item := range orderItems {
		order.Items = append(order.Items, model.OrderItem{
			Sku:   model.ProductSku(item.OrderItem.Sku),
			Count: uint16(item.OrderItem.Count),
		})
	}
	return order, nil
}

func (r *DbOrderRepository) SetStatus(ctx context.Context, orderID model.OrderID, status model.OrderStatus, inTx ...func(context.Context) error) (err error) {
	ctx, span := tracing.Start(ctx, "DbOrderRepository.SetStatus")
	defer tracing.EndWithCheckError(span, &err)

	qtx := r.queries
	if len(inTx) > 0 {
		tx, err := r.db.Begin(ctx)
		if err != nil {
			return fmt.Errorf("r.db.Begin: %w", err)
		}
		defer tx.Rollback(ctx)
		qtx = qtx.WithTx(tx)

		ctxTx := context.WithValue(ctx, repository.CtxTxKey{}, tx)
		for _, f := range inTx {
			if err := f(ctxTx); err != nil {
				return fmt.Errorf("inTx: %w", err)
			}
		}

		defer func() {
			if err = tx.Commit(ctx); err != nil {
				err = fmt.Errorf("tx.Commit: %w", err)
			}
		}()
	}

	err = qtx.SetStatus(ctx, sqlc_order.SetStatusParams{
		ID:     int64(orderID),
		Status: string(status),
	})
	if err != nil {
		return fmt.Errorf("qtx.SetStatus: %w", err)
	}
	return nil
}
