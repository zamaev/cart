package orderrepository

import (
	"context"
	"fmt"
	"route256/loms/internal/pkg/model"
)

type storage map[model.OrderID]model.Order

type orderMemoryRepository struct {
	storage storage
}

func NewOrderMemoryRepository() *orderMemoryRepository {
	return &orderMemoryRepository{
		storage: make(storage),
	}
}

func (r *orderMemoryRepository) Create(_ context.Context, order model.Order) model.OrderID {
	if order.Status == model.OrderStatusNone {
		order.Status = model.OrderStatusNew
	}
	orderID := model.OrderID(len(r.storage) + 1)
	r.storage[orderID] = order
	return orderID
}

func (r *orderMemoryRepository) GetById(_ context.Context, orderID model.OrderID) (model.Order, error) {
	order, ok := r.storage[orderID]
	if !ok {
		return model.Order{}, fmt.Errorf("invalid order id: %d", orderID)
	}
	return order, nil
}

func (r *orderMemoryRepository) SetStatus(_ context.Context, orderID model.OrderID, status model.OrderStatus) error {
	order, ok := r.storage[orderID]
	if !ok {
		return fmt.Errorf("invalid order id: %d", orderID)
	}
	order.Status = status
	r.storage[orderID] = order
	return nil
}
