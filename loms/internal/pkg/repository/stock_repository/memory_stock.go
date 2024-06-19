package stockrepository

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"route256/loms/internal/pkg/model"
)

//go:embed stock-data.json
var stockData []byte

type storage map[model.ProductSku]struct {
	TotalCount uint64 `json:"total_count"`
	Reserved   uint64 `json:"reserved"`
}

type stockMemoryRepository struct {
	storage storage
}

func NewStockMemoryRepository() (*stockMemoryRepository, error) {
	storage := make(storage)
	err := json.Unmarshal(stockData, &storage)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return &stockMemoryRepository{
		storage: storage,
	}, nil
}

func (r *stockMemoryRepository) Reserve(_ context.Context, sku model.ProductSku, count uint16) error {
	product, ok := r.storage[sku]
	if !ok {
		return fmt.Errorf("invalid sku: %d", sku)
	}
	if product.TotalCount-product.Reserved < uint64(count) {
		return fmt.Errorf("not enough products by sku: %d", sku)
	}
	product.Reserved += uint64(count)
	r.storage[sku] = product
	return nil
}

func (r *stockMemoryRepository) ReserveRemove(_ context.Context, sku model.ProductSku, count uint16) error {
	product, ok := r.storage[sku]
	if !ok {
		return fmt.Errorf("invalid sku: %d", sku)
	}
	if product.Reserved < uint64(count) {
		return fmt.Errorf("not enough reserved products by sku: %d", sku)
	}
	product.Reserved -= uint64(count)
	product.TotalCount -= uint64(count)
	r.storage[sku] = product
	return nil
}

func (r *stockMemoryRepository) ReserveCancel(_ context.Context, sku model.ProductSku, count uint16) error {
	product, ok := r.storage[sku]
	if !ok {
		return fmt.Errorf("invalid sku: %d", sku)
	}
	if product.Reserved < uint64(count) {
		return fmt.Errorf("not enough reserved products by sku: %d", sku)
	}
	product.Reserved -= uint64(count)
	r.storage[sku] = product
	return nil
}

func (r *stockMemoryRepository) GetStocksBySku(_ context.Context, sku model.ProductSku) (uint64, error) {
	product, ok := r.storage[sku]
	if !ok {
		return 0, nil
	}
	return product.TotalCount - product.Reserved, nil
}
