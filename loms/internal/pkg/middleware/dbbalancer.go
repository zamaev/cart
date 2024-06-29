package middleware

import (
	"context"
	"route256/loms/internal/pkg/utils"

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

func (db *DbBalancer) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.dbMasterPool.Exec(ctx, sql, args...)
}

func (db *DbBalancer) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	if utils.IsSqlWrite(sql) {
		return db.dbMasterPool.Query(ctx, sql, args...)
	}
	return db.GetRandomPool().Query(ctx, sql, args...)
}

func (db *DbBalancer) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	if utils.IsSqlWrite(sql) {
		return db.dbMasterPool.QueryRow(ctx, sql, args...)
	}
	return db.GetRandomPool().QueryRow(ctx, sql, args...)
}

func (db *DbBalancer) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.dbMasterPool.Begin(ctx)
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
