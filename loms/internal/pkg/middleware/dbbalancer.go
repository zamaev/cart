package middleware

import (
	"context"
	"route256/loms/internal/pkg/utils"
	"route256/loms/internal/pkg/utils/metrics"
	"route256/loms/pkg/tracing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DbBalancer struct {
	dbMasterPool  *pgxpool.Pool
	dbReplicaPool *pgxpool.Pool
}

func NewDbBalancer(dbMasterPool *pgxpool.Pool, dbReplicaPool *pgxpool.Pool) *DbBalancer {
	return &DbBalancer{
		dbMasterPool:  dbMasterPool,
		dbReplicaPool: dbReplicaPool,
	}
}

func (db *DbBalancer) Exec(ctx context.Context, sql string, args ...interface{}) (commandTag pgconn.CommandTag, err error) {
	ctx, span := tracing.Start(ctx, "DbBalancer.Exec")
	defer tracing.EndWithCheckError(span, &err)

	metrics.DbRequestCounter(utils.GetSqlType(sql))
	defer func(start time.Time) {
		metrics.DbRequestDurationWithErrorStart(utils.GetSqlType(sql), err, start)
	}(time.Now())

	commandTag, err = db.dbMasterPool.Exec(ctx, sql, args...)
	return
}

func (db *DbBalancer) Query(ctx context.Context, sql string, args ...interface{}) (rows pgx.Rows, err error) {
	ctx, span := tracing.Start(ctx, "DbBalancer.Query")
	defer tracing.EndWithCheckError(span, &err)

	metrics.DbRequestCounter(utils.GetSqlType(sql))
	defer func(start time.Time) {
		metrics.DbRequestDurationWithErrorStart(utils.GetSqlType(sql), err, start)
	}(time.Now())

	if utils.IsSqlWrite(sql) {
		rows, err = db.dbMasterPool.Query(ctx, sql, args...)
	} else {
		rows, err = db.GetRandomPool().Query(ctx, sql, args...)
	}
	return
}

func (db *DbBalancer) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	ctx, span := tracing.Start(ctx, "DbBalancer.QueryRow")
	if span != nil {
		defer span.End()
	}

	metrics.DbRequestCounter(utils.GetSqlType(sql))
	defer func(start time.Time) {
		metrics.DbRequestDurationWithErrorStart(utils.GetSqlType(sql), nil, start)
	}(time.Now())

	if utils.IsSqlWrite(sql) {
		return db.dbMasterPool.QueryRow(ctx, sql, args...)
	}
	return db.GetRandomPool().QueryRow(ctx, sql, args...)
}

func (db *DbBalancer) Begin(ctx context.Context) (_ pgx.Tx, err error) {
	ctx, span := tracing.Start(ctx, "DbBalancer.Begin")
	defer tracing.EndWithCheckError(span, &err)

	tx, err := db.dbMasterPool.Begin(ctx)
	return TxWrapper{Tx: tx}, err
}

func (db *DbBalancer) GetRandomPool() *pgxpool.Pool {
	if utils.GetRandomBool() {
		return db.dbMasterPool
	}
	return db.dbReplicaPool
}

func (db *DbBalancer) Close() {
	db.dbMasterPool.Close()
	db.dbReplicaPool.Close()
}
