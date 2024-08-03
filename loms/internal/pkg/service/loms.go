package service

import (
	"context"
	"fmt"
	"route256/loms/internal/pkg/model"
	"route256/loms/pkg/tracing"
)

type stockRepository interface {
	Reserve(context.Context, model.ProductSku, uint16) error
	ReserveRemove(context.Context, model.ProductSku, uint16) error
	ReserveCancel(context.Context, model.ProductSku, uint16) error
	GetStocksBySku(context.Context, model.ProductSku) (uint64, error)
}

type orderRepository interface {
	Create(context.Context, model.Order) (model.OrderID, error)
	GetById(context.Context, model.OrderID) (model.Order, error)
	SetStatus(context.Context, model.OrderID, model.OrderStatus) error
	GetAll(context.Context) ([]model.Order, error)
}

type LomsService struct {
	stockRepository stockRepository
	orderRepository orderRepository
}

func NewLomsService(stockRepository stockRepository, orderRepository orderRepository) *LomsService {
	return &LomsService{
		stockRepository: stockRepository,
		orderRepository: orderRepository,
	}
}

func (s *LomsService) OrderCreate(ctx context.Context, order model.Order) (_ model.OrderID, err error) {
	ctx, span := tracing.Start(ctx, "LomsService.OrderCreate")
	defer tracing.EndWithCheckError(span, &err)

	orderID, err := s.orderRepository.Create(ctx, order)
	if err != nil {
		return 0, fmt.Errorf("orderRepository.Create: %w", err)
	}
	reservedItems := make([]model.OrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		if err := s.stockRepository.Reserve(ctx, item.Sku, item.Count); err != nil {
			if err := s.orderRepository.SetStatus(ctx, orderID, model.OrderStatusFailed); err != nil {
				return 0, fmt.Errorf("orderRepository.SetStatus: %w; status: %s", err, model.OrderStatusFailed)
			}
			for i, reservedItem := range reservedItems {
				if err := s.stockRepository.ReserveCancel(ctx, reservedItem.Sku, reservedItem.Count); err != nil {
					return 0, fmt.Errorf("stockRepository.ReserveCancel: %w; not canceled reserved items: %v", err, reservedItems[i:])
				}
			}
			return 0, fmt.Errorf("stockRepository.Reserve: %w", err)
		}
		reservedItems = append(reservedItems, item)
	}
	if err := s.orderRepository.SetStatus(ctx, orderID, model.OrderStatusAwaitingPayment); err != nil {
		return 0, fmt.Errorf("orderRepository.SetStatus: %w; status: %s", err, model.OrderStatusAwaitingPayment)
	}
	return orderID, nil
}

func (s *LomsService) OrderInfo(ctx context.Context, orderID model.OrderID) (_ model.Order, err error) {
	ctx, span := tracing.Start(ctx, "LomsService.OrderInfo")
	defer tracing.EndWithCheckError(span, &err)

	return s.orderRepository.GetById(ctx, orderID)
}

func (s *LomsService) OrderPay(ctx context.Context, orderID model.OrderID) (err error) {
	ctx, span := tracing.Start(ctx, "LomsService.OrderPay")
	defer tracing.EndWithCheckError(span, &err)

	order, err := s.orderRepository.GetById(ctx, orderID)
	if err != nil {
		return fmt.Errorf("orderRepository.GetById: %w", err)
	}
	if order.Status != model.OrderStatusAwaitingPayment {
		return fmt.Errorf("invalid order status: %s; orderID: %d", order.Status, orderID)
	}
	for i, item := range order.Items {
		if err := s.stockRepository.ReserveRemove(ctx, item.Sku, item.Count); err != nil {
			return fmt.Errorf("stockRepository.ReserveRemove: %w; not removed items: %v", err, order.Items[i:])
		}
	}
	return s.orderRepository.SetStatus(ctx, orderID, model.OrderStatusPaid)
}

func (s *LomsService) OrderCancel(ctx context.Context, orderID model.OrderID) (err error) {
	ctx, span := tracing.Start(ctx, "LomsService.OrderCancel")
	defer tracing.EndWithCheckError(span, &err)

	order, err := s.orderRepository.GetById(ctx, orderID)
	if err != nil {
		return fmt.Errorf("orderRepository.GetById: %w", err)
	}
	if order.Status != model.OrderStatusAwaitingPayment {
		return fmt.Errorf("invalid order status: %s; orderID: %d", order.Status, orderID)
	}
	for i, item := range order.Items {
		if err := s.stockRepository.ReserveCancel(ctx, item.Sku, item.Count); err != nil {
			return fmt.Errorf("stockRepository.ReserveCancel: %w; not canceled items: %v", err, order.Items[i:])
		}
	}
	return s.orderRepository.SetStatus(ctx, orderID, model.OrderStatusCancelled)
}

func (s *LomsService) StocksInfo(ctx context.Context, sku model.ProductSku) (_ uint64, err error) {
	ctx, span := tracing.Start(ctx, "LomsService.StocksInfo")
	defer tracing.EndWithCheckError(span, &err)

	return s.stockRepository.GetStocksBySku(ctx, sku)
}

func (s *LomsService) GetAllOrders(ctx context.Context) (_ []model.Order, err error) {
	ctx, span := tracing.Start(ctx, "GetAllOrders")
	defer tracing.EndWithCheckError(span, &err)

	return s.orderRepository.GetAll(ctx)
}
