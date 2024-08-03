package orderrepository

import (
	"context"
	"fmt"
	"route256/loms/internal/pkg/inrfa/shard_manager"
	"route256/loms/internal/pkg/model"
	"route256/loms/internal/pkg/repository"
	"route256/loms/internal/pkg/repository/order_repository/sqlc_order"
	"route256/loms/pkg/tracing"
	"sort"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type shardManager interface {
	GetShardIndex(key shard_manager.ShardKey) shard_manager.ShardIndex
	GetShardIndexFromID(id int64) shard_manager.ShardIndex
	Pick(key shard_manager.ShardIndex) (*pgxpool.Pool, error)
	GetShards() []*pgxpool.Pool
}

type DbOrderRepository struct {
	sm shardManager
}

func NewDbOrderRepository(sm shardManager) *DbOrderRepository {
	return &DbOrderRepository{
		sm: sm,
	}
}

func (r *DbOrderRepository) Create(ctx context.Context, order model.Order, inTx ...func(context.Context, model.OrderID, model.OrderStatus) error) (_ model.OrderID, err error) {
	ctx, span := tracing.Start(ctx, "DbOrderRepository.Create")
	defer tracing.EndWithCheckError(span, &err)

	shIndex := r.sm.GetShardIndex(shard_manager.ShardKey(strconv.FormatInt(int64(order.User), 10)))
	db, err := r.sm.Pick(shIndex)
	if err != nil {
		return 0, fmt.Errorf("r.sm.Pick: %w", err)
	}
	queries := sqlc_order.New(db)

	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("r.db.Begin: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := queries.WithTx(tx)

	if order.Status == model.OrderStatusNone {
		order.Status = model.OrderStatusNew
	}

	id, err := qtx.Create(ctx, sqlc_order.CreateParams{
		UserID:  int64(order.User),
		Status:  string(order.Status),
		ShardID: int32(shIndex),
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

	shIndex := r.sm.GetShardIndexFromID(int64(orderID))
	db, err := r.sm.Pick(shIndex)
	if err != nil {
		return model.Order{}, fmt.Errorf("r.sm.Pick: %w", err)
	}
	queries := sqlc_order.New(db)

	orderItems, err := queries.GetById(ctx, int64(orderID))
	if err != nil {
		return model.Order{}, fmt.Errorf("r.queries.GetById: %w", err)
	}
	if len(orderItems) == 0 {
		return model.Order{}, fmt.Errorf("invalid order id: %d", orderID)
	}
	order := model.Order{
		ID:     model.OrderID(orderItems[0].Order.ID),
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

	shIndex := r.sm.GetShardIndexFromID(int64(orderID))
	db, err := r.sm.Pick(shIndex)
	if err != nil {
		return fmt.Errorf("r.sm.Pick: %w", err)
	}
	queries := sqlc_order.New(db)

	qtx := queries
	if len(inTx) > 0 {
		tx, err := db.Begin(ctx)
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

func (r *DbOrderRepository) GetAll(ctx context.Context) (_ []model.Order, err error) {
	ctx, span := tracing.Start(ctx, "DbOrderRepository.GetAll")
	defer tracing.EndWithCheckError(span, &err)

	orders := make([]model.Order, 0, 10)

	for _, db := range r.sm.GetShards() {
		queries := sqlc_order.New(db)
		orderItems, err := queries.GetAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("r.queries.GetAll: %w", err)
		}
		if len(orderItems) == 0 {
			continue
		}

		var order model.Order
		for _, item := range orderItems {
			// Начался новый заказ
			if item.Order.ID != int64(order.ID) {
				// Старый сохраняем
				if order.ID != 0 {
					orders = append(orders, order)
				}
				order = model.Order{
					ID:     model.OrderID(item.Order.ID),
					Status: model.OrderStatus(item.Order.Status),
					User:   model.UserID(item.Order.UserID),
					Items:  make([]model.OrderItem, 0, 3),
				}
			}
			order.Items = append(order.Items, model.OrderItem{
				Sku:   model.ProductSku(item.OrderItem.Sku),
				Count: uint16(item.OrderItem.Count),
			})

		}
		orders = append(orders, order)
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].ID > orders[j].ID
	})

	return orders, nil
}
