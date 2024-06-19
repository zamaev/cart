package server

import (
	"context"
	"fmt"
	"log"
	"route256/loms/internal/pkg/model"
	orderrepository "route256/loms/internal/pkg/repository/order_repository"
	stockrepository "route256/loms/internal/pkg/repository/stock_repository"
	"route256/loms/internal/pkg/service"
	"route256/loms/pkg/api/loms/v1"
)

type LomsService interface {
	OrderCreate(ctx context.Context, order model.Order) (model.OrderID, error)
	OrderInfo(ctx context.Context, orderID model.OrderID) (model.Order, error)
	OrderPay(ctx context.Context, orderID model.OrderID) error
	OrderCancel(ctx context.Context, orderID model.OrderID) error
	StocksInfo(ctx context.Context, sku model.ProductSku) (uint64, error)
}

type Server struct {
	loms.UnimplementedLomsServer
	service LomsService
}

func NewServer() *Server {
	stockRepository, err := stockrepository.NewStockMemoryRepository()
	if err != nil {
		log.Fatal(err)
	}
	orderRepository := orderrepository.NewOrderMemoryRepository()
	lomsService := service.NewLomsService(stockRepository, orderRepository)

	return &Server{
		service: lomsService,
	}
}

func (s *Server) OrderCreate(ctx context.Context, req *loms.OrderCreateRequest) (*loms.OrderCreateResponse, error) {
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
		return nil, fmt.Errorf("lomsService.OrderCreate: %w", err)
	}
	return &loms.OrderCreateResponse{
		OrderId: int64(orderId),
	}, nil
}

func (s *Server) OrderInfo(ctx context.Context, req *loms.OrderInfoRequest) (*loms.OrderInfoResponse, error) {
	order, err := s.service.OrderInfo(ctx, model.OrderID(req.OrderId))
	if err != nil {
		return nil, fmt.Errorf("lomsService.OrderInfo: %w", err)
	}
	res := &loms.OrderInfoResponse{
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

func (s *Server) OrderPay(ctx context.Context, req *loms.OrderPayRequest) (*loms.OrderPayResponse, error) {
	err := s.service.OrderPay(ctx, model.OrderID(req.OrderId))
	if err != nil {
		return nil, fmt.Errorf("lomsService.OrderPay: %w", err)
	}
	return nil, nil
}
func (s *Server) OrderCancel(ctx context.Context, req *loms.OrderCancelRequest) (*loms.OrderCancelResponse, error) {
	err := s.service.OrderCancel(ctx, model.OrderID(req.OrderId))
	if err != nil {
		return nil, fmt.Errorf("lomsService.OrderCancel: %w", err)
	}
	return nil, nil
}
func (s *Server) StocksInfo(ctx context.Context, req *loms.StocksInfoRequest) (*loms.StocksInfoResponse, error) {
	count, err := s.service.StocksInfo(ctx, model.ProductSku(req.Sku))
	if err != nil {
		return nil, fmt.Errorf("lomsService.StocksInfo: %w", err)
	}
	return &loms.StocksInfoResponse{
		Count: uint64(count),
	}, nil
}
