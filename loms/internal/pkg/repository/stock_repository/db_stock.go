package stockrepository

import (
	"context"
	"fmt"
	"route256/loms/internal/pkg/model"
	"route256/loms/internal/pkg/repository/stock_repository/sqlc_stock"
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

type DbStockRepository struct {
	queries *sqlc_stock.Queries
}

func NewDbStockRepository(db DB) (*DbStockRepository, error) {
	return &DbStockRepository{
		queries: sqlc_stock.New(db),
	}, nil
}

func (r *DbStockRepository) Reserve(ctx context.Context, sku model.ProductSku, count uint16) (err error) {
	ctx, span := tracing.Start(ctx, "DbStockRepository.Reserve")
	defer tracing.EndWithCheckError(span, &err)

	stockCount, err := r.GetStocksBySku(ctx, sku)
	if err != nil {
		return fmt.Errorf("r.GetStocksBySku: %w", err)
	}
	if stockCount < uint64(count) {
		return fmt.Errorf("not enough products in stock")
	}
	return r.queries.Reserve(ctx, sqlc_stock.ReserveParams{
		Sku:   int64(sku),
		Count: int32(count),
	})
}

func (r *DbStockRepository) ReserveRemove(ctx context.Context, sku model.ProductSku, count uint16) (err error) {
	ctx, span := tracing.Start(ctx, "DbStockRepository.ReserveRemove")
	defer tracing.EndWithCheckError(span, &err)

	return r.queries.ReserveRemove(ctx, sqlc_stock.ReserveRemoveParams{
		Sku:   int64(sku),
		Count: int32(count),
	})
}

func (r *DbStockRepository) ReserveCancel(ctx context.Context, sku model.ProductSku, count uint16) (err error) {
	ctx, span := tracing.Start(ctx, "DbStockRepository.ReserveCancel")
	defer tracing.EndWithCheckError(span, &err)

	return r.queries.ReserveCancel(ctx, sqlc_stock.ReserveCancelParams{
		Sku:   int64(sku),
		Count: int32(count),
	})
}

func (r *DbStockRepository) GetStocksBySku(ctx context.Context, sku model.ProductSku) (_ uint64, err error) {
	ctx, span := tracing.Start(ctx, "DbStockRepository.GetStocksBySku")
	defer tracing.EndWithCheckError(span, &err)

	count, err := r.queries.GetStocksBySku(ctx, int64(sku))
	if err != nil {
		return 0, fmt.Errorf("r.queries.GetStocksBySku: %w", err)
	}
	return uint64(count), nil
}
