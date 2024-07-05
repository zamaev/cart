package repository

import (
	"context"
	"sync"
	"testing"
)

const count = 100000

func TestRaceAddProduct(t *testing.T) {
	t.Parallel()

	repo := NewCartMemoryRepository()

	wg := &sync.WaitGroup{}
	for range count {
		wg.Add(1)
		go func() {
			repo.AddProduct(context.Background(), 1, 1, 1)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestRaceRemoveProduct(t *testing.T) {
	t.Parallel()

	repo := NewCartMemoryRepository()

	wg := &sync.WaitGroup{}
	for range count {
		wg.Add(1)
		go func() {
			repo.RemoveProduct(context.Background(), 1, 1)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestRaceClearCart(t *testing.T) {
	t.Parallel()

	repo := NewCartMemoryRepository()

	wg := &sync.WaitGroup{}
	for range count {
		wg.Add(1)
		go func() {
			repo.ClearCart(context.Background(), 1)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestRaceGetCart(t *testing.T) {
	t.Parallel()

	repo := NewCartMemoryRepository()

	wg := &sync.WaitGroup{}
	for range count {
		wg.Add(1)
		go func() {
			repo.GetCart(context.Background(), 1)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestRaceGetProductCount(t *testing.T) {
	t.Parallel()

	repo := NewCartMemoryRepository()

	wg := &sync.WaitGroup{}
	for range count {
		wg.Add(1)
		go func() {
			repo.GetProductCount(context.Background(), 1, 1)
			wg.Done()
		}()
	}
	wg.Wait()
}
