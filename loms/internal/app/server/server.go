package server

import (
	"context"
	"fmt"
	"route256/loms/internal/pkg/config"
	"route256/loms/internal/pkg/inrfa/kafka/producer"
	"route256/loms/internal/pkg/middleware"
	"route256/loms/internal/pkg/model"
	orderrepository "route256/loms/internal/pkg/repository/order_repository"
	outboxrepository "route256/loms/internal/pkg/repository/outbox_repository"
	stockrepository "route256/loms/internal/pkg/repository/stock_repository"
	"route256/loms/internal/pkg/service"
	"route256/loms/pkg/api/loms/v1"
	"route256/loms/pkg/logger"
	"route256/loms/pkg/tracing"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LomsService interface {
	OrderCreate(ctx context.Context, order model.Order) (model.OrderID, error)
	OrderInfo(ctx context.Context, orderID model.OrderID) (model.Order, error)
	OrderPay(ctx context.Context, orderID model.OrderID) error
	OrderCancel(ctx context.Context, orderID model.OrderID) error
	StocksInfo(ctx context.Context, sku model.ProductSku) (uint64, error)
}

type Producer interface {
	Close()
	RunEventsHandle()
}

type Server struct {
	loms.UnimplementedLomsServer
	service  LomsService
	producer Producer
}

func NewServer(config config.Config) *Server {
	ctx := context.Background()

	dbMasterPool, err := pgxpool.New(ctx, config.DbMasterUrl)
	if err != nil {
		logger.Panicw(ctx, "create connection master pool", "err", err)
	}
	if err := dbMasterPool.Ping(ctx); err != nil {
		logger.Panicw(ctx, "dbMasterPool.Ping", "err", err)
	}
	dbReplicaPool, err := pgxpool.New(ctx, config.DbReplicaUrl)
	if err != nil {
		logger.Panicw(ctx, "create connection replica pool", "err", err)
	}
	if err := dbReplicaPool.Ping(ctx); err != nil {
		logger.Panicw(ctx, "dbReplicaPool.Ping", "err", err)
	}
	dbBalancer := middleware.NewDbBalancer(dbMasterPool, dbReplicaPool)

	stockRepository, err := stockrepository.NewDbStockRepository(dbBalancer)
	if err != nil {
		logger.Panicw(ctx, "stockrepository.NewDbStockRepository", "err", err)
	}
	orderRepository := orderrepository.NewDbOrderRepository(dbBalancer)
	outboxRepository := outboxrepository.NewDbOutboxRepository(dbBalancer)
	orderOutboxWrapper := middleware.NewOrderOutboxWrapper(orderRepository, outboxRepository, config.Kafka)
	lomsService := service.NewLomsService(stockRepository, orderOutboxWrapper)

	prod, err := producer.NewProducer(ctx, config.Kafka, outboxRepository,
		producer.WithProducerPartitioner(sarama.NewHashPartitioner),
		producer.WithRequiredAcks(sarama.WaitForAll),
		producer.WithMaxRetries(5),
		producer.WithRetryBackoff(10*time.Millisecond),
		producer.WithMaxOpenRequests(1),
		producer.WithProducerCompression(sarama.CompressionGZIP),
		producer.WithProducerReturnSuccesse(),
		producer.WithProducerFlushMessages(8),
		producer.WithProducerFlushFrequency(5*time.Second),
	)
	if err != nil {
		logger.Panicw(ctx, "producer.NewProducer", "err", err)
	}
	go prod.RunEventsHandle()

	return &Server{
		service:  lomsService,
		producer: prod,
	}
}

func (s *Server) Shutdown() {
	s.producer.Close()
}

func (s *Server) OrderCreate(ctx context.Context, req *loms.OrderCreateRequest) (res *loms.OrderCreateResponse, err error) {
	ctx, span := tracing.Start(ctx, "Server.OrderCreate")
	defer tracing.EndWithCheckError(span, &err)

	order := model.Order{
		User:  model.UserID(req.User),
		Items: make([]model.OrderItem, 0, len(req.Items)),
	}
	for _, item := range req.Items {
		order.Items = append(order.Items, model.OrderItem{
			Sku:   model.ProductSku(item.Sku),
			Count: uint16(item.Count),
		})
	}
	orderId, err := s.service.OrderCreate(ctx, order)
	if err != nil {
		logger.Errorw(ctx, "lomsService.OrderCreate", "err", err)
		return nil, fmt.Errorf("lomsService.OrderCreate: %w", err)
	}
	return &loms.OrderCreateResponse{
		OrderId: int64(orderId),
	}, nil
}

func (s *Server) OrderInfo(ctx context.Context, req *loms.OrderInfoRequest) (res *loms.OrderInfoResponse, err error) {
	ctx, span := tracing.Start(ctx, "Server.OrderInfo")
	defer tracing.EndWithCheckError(span, &err)

	order, err := s.service.OrderInfo(ctx, model.OrderID(req.OrderId))
	if err != nil {
		logger.Errorw(ctx, "lomsService.OrderInfo", "err", err)
		return nil, fmt.Errorf("lomsService.OrderInfo: %w", err)
	}
	res = &loms.OrderInfoResponse{
		Status: string(order.Status),
		User:   int64(order.User),
		Items:  make([]*loms.OrderItem, 0, len(order.Items)),
	}
	for _, item := range order.Items {
		res.Items = append(res.Items, &loms.OrderItem{
			Sku:   uint32(item.Sku),
			Count: uint32(item.Count),
		})
	}
	return res, nil
}

func (s *Server) OrderPay(ctx context.Context, req *loms.OrderPayRequest) (res *loms.OrderPayResponse, err error) {
	ctx, span := tracing.Start(ctx, "Server.OrderPay")
	defer tracing.EndWithCheckError(span, &err)

	err = s.service.OrderPay(ctx, model.OrderID(req.OrderId))
	if err != nil {
		logger.Errorw(ctx, "lomsService.OrderPay", "err", err)
		return nil, fmt.Errorf("lomsService.OrderPay: %w", err)
	}
	return nil, nil
}

func (s *Server) OrderCancel(ctx context.Context, req *loms.OrderCancelRequest) (res *loms.OrderCancelResponse, err error) {
	ctx, span := tracing.Start(ctx, "Server.OrderCancel")
	defer tracing.EndWithCheckError(span, &err)

	err = s.service.OrderCancel(ctx, model.OrderID(req.OrderId))
	if err != nil {
		logger.Errorw(ctx, "lomsService.OrderCancel", "err", err)
		return nil, fmt.Errorf("lomsService.OrderCancel: %w", err)
	}
	return nil, nil
}

func (s *Server) StocksInfo(ctx context.Context, req *loms.StocksInfoRequest) (res *loms.StocksInfoResponse, err error) {
	ctx, span := tracing.Start(ctx, "Server.StocksInfo")
	defer tracing.EndWithCheckError(span, &err)

	count, err := s.service.StocksInfo(ctx, model.ProductSku(req.Sku))
	if err != nil {
		logger.Errorw(ctx, "lomsService.StocksInfo", "err", err)
		return nil, fmt.Errorf("lomsService.StocksInfo: %w", err)
	}
	return &loms.StocksInfoResponse{
		Count: uint64(count),
	}, nil
}
