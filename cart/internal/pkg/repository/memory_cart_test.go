package repository

import (
	"context"
	"route256/cart/internal/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddProduct(t *testing.T) {
	ctx := context.Background()
	repo := NewCartMemoryRepository()

	err := repo.AddProduct(ctx, 1, 1, 1)
	assert.NoError(t, err)
	assert.Len(t, repo.storage, 1)

	err = repo.AddProduct(ctx, 2, 1, 1)
	assert.NoError(t, err)
	assert.Len(t, repo.storage, 2)
}

func TestRemoveProduct(t *testing.T) {
	ctx := context.Background()
	repo := NewCartMemoryRepository()

	repo.AddProduct(ctx, 1, 1, 1)
	assert.Len(t, repo.storage[1], 1)

	err := repo.RemoveProduct(ctx, 1, 1)
	assert.NoError(t, err)
	assert.Len(t, repo.storage[1], 0)
}

func TestClearCart(t *testing.T) {
	ctx := context.Background()
	repo := NewCartMemoryRepository()

	repo.AddProduct(ctx, 1, 1, 1)
	repo.AddProduct(ctx, 2, 1, 1)
	assert.Len(t, repo.storage, 2)

	err := repo.ClearCart(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, repo.storage[1], 0)
}

func TestGetCart(t *testing.T) {
	ctx := context.Background()
	repo := NewCartMemoryRepository()

	repo.AddProduct(ctx, 1, 1, 5)
	repo.AddProduct(ctx, 1, 2, 7)

	cart, err := repo.GetCart(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, cart, model.Cart{
		1: 5,
		2: 7,
	})
}

func BenchmarkAddProduct(b *testing.B) {
	ctx := context.Background()
	repo := NewCartMemoryRepository()
	for i := 0; i < b.N; i++ {
		repo.AddProduct(ctx, 1, 1, 1)
	}
}
