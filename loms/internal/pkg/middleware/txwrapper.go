package middleware

import (
	"context"
	"route256/loms/internal/pkg/utils"
	"route256/loms/internal/pkg/utils/metrics"
	"route256/loms/pkg/tracing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type TxWrapper struct {
	pgx.Tx
}

func (tx TxWrapper) Exec(ctx context.Context, sql string, args ...interface{}) (commandTag pgconn.CommandTag, err error) {
	ctx, span := tracing.Start(ctx, "TxWrapper.Exec")
	defer tracing.EndWithCheckError(span, &err)

	metrics.DbRequestCounter(utils.GetSqlType(sql))
	defer func(start time.Time) {
		metrics.DbRequestDurationWithErrorStart(utils.GetSqlType(sql), err, start)
	}(time.Now())

	return tx.Tx.Exec(ctx, sql, args...)
}
func (tx TxWrapper) Query(ctx context.Context, sql string, args ...interface{}) (rows pgx.Rows, err error) {
	ctx, span := tracing.Start(ctx, "TxWrapper.Query")
	defer tracing.EndWithCheckError(span, &err)

	metrics.DbRequestCounter(utils.GetSqlType(sql))
	defer func(start time.Time) {
		metrics.DbRequestDurationWithErrorStart(utils.GetSqlType(sql), err, start)
	}(time.Now())

	return tx.Tx.Query(ctx, sql, args...)
}
func (tx TxWrapper) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	ctx, span := tracing.Start(ctx, "TxWrapper.QueryRow")
	defer span.End()

	metrics.DbRequestCounter(utils.GetSqlType(sql))
	defer func(start time.Time) {
		metrics.DbRequestDurationWithErrorStart(utils.GetSqlType(sql), nil, start)
	}(time.Now())

	return tx.Tx.QueryRow(ctx, sql, args...)
}

func (tx TxWrapper) Commit(ctx context.Context) (err error) {
	ctx, span := tracing.Start(ctx, "TxWrapper.Commit")
	defer tracing.EndWithCheckError(span, &err)

	return tx.Tx.Commit(ctx)
}
