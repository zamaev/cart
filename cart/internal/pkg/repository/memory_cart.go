package repository

import (
	"context"
	"route256/cart/internal/pkg/model"
)

type Storage map[model.UserId]model.Cart

type CartMemoryRepository struct {
	storage Storage
}

func NewCartMemoryRepository() *CartMemoryRepository {
	return &CartMemoryRepository{
		storage: make(Storage),
	}
}

func (r *CartMemoryRepository) AddProduct(_ context.Context, userId model.UserId, ProductSku model.ProductSku, count uint16) error {
	if _, ok := r.storage[userId]; !ok {
		r.storage[userId] = make(model.Cart)
	}
	r.storage[userId][ProductSku] += count
	return nil
}

func (r *CartMemoryRepository) RemoveProduct(_ context.Context, userId model.UserId, ProductSku model.ProductSku) error {
	if _, ok := r.storage[userId]; !ok {
		return nil
	}
	delete(r.storage[userId], ProductSku)
	return nil
}

func (r *CartMemoryRepository) ClearCart(_ context.Context, userId model.UserId) error {
	delete(r.storage, userId)
	return nil
}

func (r *CartMemoryRepository) GetCart(_ context.Context, userId model.UserId) (model.Cart, error) {
	if _, ok := r.storage[userId]; !ok {
		r.storage[userId] = make(model.Cart)
	}
	return r.storage[userId], nil
}
