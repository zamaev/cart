package test

import (
	"context"
	"database/sql"
	"route256/loms/internal/pkg/config"
	"route256/loms/internal/pkg/middleware"
	"route256/loms/test/testconfig"
	"testing"
	"time"

	stockrepository "route256/loms/internal/pkg/repository/stock_repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const stockMigrationVersion = 20240620195318

type StockSute struct {
	suite.Suite
	db               *sql.DB
	stockRepository  *stockrepository.DbStockRepository
	migrationsDownTo int64
}

func TestStockSuite(t *testing.T) {
	suite.Run(t, new(StockSute))
}

func (s *StockSute) SetupSuite() {
	config := config.NewConfig()
	testconfig := testconfig.NewConfig()

	dbMasterPool, err := pgxpool.New(context.Background(), config.DbMasterUrl)
	require.NoError(s.T(), err)
	err = dbMasterPool.Ping(context.Background())
	require.NoError(s.T(), err)

	dbReplicaPool, err := pgxpool.New(context.Background(), config.DbReplicaUrl)
	require.NoError(s.T(), err)
	err = dbReplicaPool.Ping(context.Background())
	require.NoError(s.T(), err)

	dbBalancer := middleware.NewDbBalancer(dbMasterPool, dbReplicaPool)

	db := stdlib.OpenDBFromPool(dbMasterPool)
	s.db = db

	err = goose.SetDialect("postgres")
	require.NoError(s.T(), err)

	version, err := goose.GetDBVersion(db)
	require.NoError(s.T(), err)
	if version >= stockMigrationVersion && !testconfig.CheckMigratedTests {
		s.T().Skip()
	}
	require.Less(s.T(), version, int64(stockMigrationVersion), "db version should be less than %d", stockMigrationVersion)

	s.migrationsDownTo = version

	err = goose.UpTo(db, "../migrations", stockMigrationVersion)
	require.NoError(s.T(), err)

	// Waiting for data to reach replication
	time.Sleep(1 * time.Second)

	stockRepository, err := stockrepository.NewDbStockRepository(dbBalancer)
	require.NoError(s.T(), err)

	s.stockRepository = stockRepository
}

func (s *StockSute) TearDownSuite() {
	err := goose.DownTo(s.db, "../migrations", s.migrationsDownTo)
	require.NoError(s.T(), err)
}

// Waiting for data to reach replication
func (s *StockSute) SetupTest() {
	time.Sleep(time.Millisecond)
}

func (s *StockSute) TestAGetStocksBySku() {
	ctx := context.Background()
	count, err := s.stockRepository.GetStocksBySku(ctx, 1)
	require.NoError(s.T(), err)
	require.Equal(s.T(), uint64(4), count)
}

func (s *StockSute) TestBReverse() {
	ctx := context.Background()
	err := s.stockRepository.Reserve(ctx, 1, 1)
	require.NoError(s.T(), err)
}

func (s *StockSute) TestCGetStocksBySkuAfterReverse() {
	ctx := context.Background()
	count, err := s.stockRepository.GetStocksBySku(ctx, 1)
	require.NoError(s.T(), err)
	require.Equal(s.T(), uint64(3), count)
}

func (s *StockSute) TestDReverseCancel() {
	ctx := context.Background()
	err := s.stockRepository.ReserveCancel(ctx, 1, 1)
	require.NoError(s.T(), err)
}

func (s *StockSute) TestEGetStocksBySkuAfterReverseCancel() {
	ctx := context.Background()
	count, err := s.stockRepository.GetStocksBySku(ctx, 1)
	require.NoError(s.T(), err)
	require.Equal(s.T(), uint64(4), count)
}

func (s *StockSute) TestFReverse() {
	ctx := context.Background()
	err := s.stockRepository.Reserve(ctx, 1, 1)
	require.NoError(s.T(), err)
}

func (s *StockSute) TestGGetStocksBySkuAfterReverse() {
	ctx := context.Background()
	count, err := s.stockRepository.GetStocksBySku(ctx, 1)
	require.NoError(s.T(), err)
	require.Equal(s.T(), uint64(3), count)
}

func (s *StockSute) TestHReverseRemove() {
	ctx := context.Background()
	err := s.stockRepository.ReserveRemove(ctx, 1, 1)
	require.NoError(s.T(), err)
}

func (s *StockSute) TestIGetStocksBySkuAfterReverseRemove() {
	ctx := context.Background()
	count, err := s.stockRepository.GetStocksBySku(ctx, 1)
	require.NoError(s.T(), err)
	require.Equal(s.T(), uint64(3), count)
}

func (s *StockSute) TestJGetStocksBySkuError() {
	ctx := context.Background()
	_, err := s.stockRepository.GetStocksBySku(ctx, 888888)
	require.Error(s.T(), err)
}
