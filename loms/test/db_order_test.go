package test

import (
	"context"
	"database/sql"
	"route256/loms/internal/pkg/config"
	"route256/loms/internal/pkg/inrfa/shard_manager"
	"route256/loms/internal/pkg/model"
	"route256/loms/test/testconfig"
	"testing"
	"time"

	orderrepository "route256/loms/internal/pkg/repository/order_repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const orderMigrationVersion = 20240620195339

type OrderSute struct {
	suite.Suite
	db               *sql.DB
	orderRepository  *orderrepository.DbOrderRepository
	migrationsDownTo int64
}

func TestOrderSute(t *testing.T) {
	suite.Run(t, new(OrderSute))
}

func (s *OrderSute) SetupSuite() {
	ctx := context.Background()

	config, err := config.NewConfig()
	require.NoError(s.T(), err)
	testconfig := testconfig.NewConfig()

	dbMasterPool, err := pgxpool.New(ctx, config.DbMasterUrl)
	require.NoError(s.T(), err)
	err = dbMasterPool.Ping(ctx)
	require.NoError(s.T(), err)

	dbMasterShard2Pool, err := pgxpool.New(ctx, config.DbMasterShard2Url)
	require.NoError(s.T(), err)
	err = dbMasterShard2Pool.Ping(ctx)
	require.NoError(s.T(), err)
	shardManager := shard_manager.NewShardManager(
		shard_manager.GetMurmur3ShardFn(2),
		dbMasterPool,
		dbMasterShard2Pool,
	)

	db := stdlib.OpenDBFromPool(dbMasterPool)
	s.db = db

	err = goose.SetDialect("postgres")
	require.NoError(s.T(), err)

	version, err := goose.GetDBVersion(db)
	require.NoError(s.T(), err)
	if version >= orderMigrationVersion && !testconfig.CheckMigratedTests {
		s.T().Skip()
	}
	require.Less(s.T(), version, int64(orderMigrationVersion), "db version should be less than %d", orderMigrationVersion)

	s.migrationsDownTo = version

	err = goose.UpTo(db, "../migrations", orderMigrationVersion)
	require.NoError(s.T(), err)

	// Waiting for data to reach replication
	time.Sleep(1 * time.Second)

	orderRepository := orderrepository.NewDbOrderRepository(shardManager)

	s.orderRepository = orderRepository
}

func (s *OrderSute) TearDownSuite() {
	err := goose.DownTo(s.db, "../migrations", s.migrationsDownTo)
	require.NoError(s.T(), err)
}

// Waiting for data to reach replication
func (s *OrderSute) SetupTest() {
	time.Sleep(time.Millisecond)
}

func (s *OrderSute) TestACreateOrder() {
	ctx := context.Background()
	order := model.Order{
		Status: model.OrderStatusNew,
		User:   1,
		Items: []model.OrderItem{
			{
				Sku:   1,
				Count: 1,
			},
		},
	}
	orderId, err := s.orderRepository.Create(ctx, order)
	require.NoError(s.T(), err)
	require.Equal(s.T(), model.OrderID(1), orderId)
}

func (s *OrderSute) TestBGetOrder() {
	ctx := context.Background()
	order, err := s.orderRepository.GetById(ctx, 1)
	require.NoError(s.T(), err)
	require.Equal(s.T(), model.Order{
		Status: model.OrderStatusNew,
		User:   1,
		Items: []model.OrderItem{
			{
				Sku:   1,
				Count: 1,
			},
		},
	}, order)
}

func (s *OrderSute) TestCSetOrderStatus() {
	ctx := context.Background()
	err := s.orderRepository.SetStatus(ctx, 1, model.OrderStatusPaid)
	require.NoError(s.T(), err)
}

func (s *OrderSute) TestDCheckOrderStatus() {
	ctx := context.Background()
	order, err := s.orderRepository.GetById(ctx, 1)
	require.NoError(s.T(), err)
	require.Equal(s.T(), model.OrderStatusPaid, order.Status)
}

func (s *OrderSute) TestEGetOrderError() {
	ctx := context.Background()
	_, err := s.orderRepository.GetById(ctx, 2)
	require.Error(s.T(), err)
}
