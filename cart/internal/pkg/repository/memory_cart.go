package repository

import (
	"context"
	"route256/cart/internal/pkg/model"
	"route256/cart/internal/pkg/utils/metrics"
	"route256/cart/pkg/tracing"
	"sync"
)

type Storage map[model.UserId]model.Cart

type CartMemoryRepository struct {
	mx      sync.RWMutex
	storage Storage
}

func NewCartMemoryRepository() *CartMemoryRepository {
	return &CartMemoryRepository{
		storage: make(Storage),
	}
}

func (r *CartMemoryRepository) AddProduct(ctx context.Context, userId model.UserId, ProductSku model.ProductSku, count uint16) error {
	_, span := tracing.Start(ctx, "CartMemoryRepository.AddProduct")
	defer span.End()

	r.mx.Lock()
	defer r.mx.Unlock()
	if _, ok := r.storage[userId]; !ok {
		r.storage[userId] = make(model.Cart)
	}
	r.storage[userId][ProductSku] += count
	r.SendMetrics()
	return nil
}

func (r *CartMemoryRepository) RemoveProduct(ctx context.Context, userId model.UserId, ProductSku model.ProductSku) error {
	_, span := tracing.Start(ctx, "CartMemoryRepository.RemoveProduct")
	defer span.End()

	r.mx.Lock()
	defer r.mx.Unlock()
	if _, ok := r.storage[userId]; !ok {
		return nil
	}
	delete(r.storage[userId], ProductSku)
	r.SendMetrics()
	return nil
}

func (r *CartMemoryRepository) ClearCart(ctx context.Context, userId model.UserId) error {
	_, span := tracing.Start(ctx, "CartMemoryRepository.ClearCart")
	defer span.End()

	r.mx.Lock()
	defer r.mx.Unlock()
	delete(r.storage, userId)
	r.SendMetrics()
	return nil
}

func (r *CartMemoryRepository) GetCart(ctx context.Context, userId model.UserId) (model.Cart, error) {
	_, span := tracing.Start(ctx, "CartMemoryRepository.GetCart")
	defer span.End()

	r.mx.Lock()
	defer r.mx.Unlock()
	if _, ok := r.storage[userId]; !ok {
		r.storage[userId] = make(model.Cart)
	}
	r.SendMetrics()
	return r.storage[userId], nil
}

func (r *CartMemoryRepository) GetProductCount(ctx context.Context, userId model.UserId, ProductSku model.ProductSku) (uint16, error) {
	_, span := tracing.Start(ctx, "CartMemoryRepository.GetProductCount")
	defer span.End()

	r.mx.RLock()
	defer r.mx.RUnlock()
	if _, ok := r.storage[userId]; !ok {
		return 0, nil
	}
	r.SendMetrics()
	return r.storage[userId][ProductSku], nil
}

func (r *CartMemoryRepository) SendMetrics() {
	var amount float64
	for _, cart := range r.storage {
		for _, count := range cart {
			amount += float64(count)
		}
	}
	metrics.CartRepositoryAmounter(amount)
}
