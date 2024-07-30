package product_cache

import (
	"context"
	"errors"
	"route256/cart/internal/pkg/model"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type ProductServiceTest struct {
	counter int
	mx      sync.Mutex
}

func (p *ProductServiceTest) GetProduct(_ context.Context, _ model.ProductSku) (*model.Product, error) {
	time.Sleep(1 * time.Second)
	p.mx.Lock()
	defer p.mx.Unlock()
	p.counter++
	return &model.Product{
		Sku:  1,
		Name: "test",
	}, nil
}

type CacheTest struct {
	cache    map[string]string
	mx       sync.Mutex
	getCount int
	setCount int
}

func (c *CacheTest) Get(ctx context.Context, key string) (string, error) {
	c.mx.Lock()
	c.getCount++
	defer c.mx.Unlock()
	if val, ok := c.cache[key]; ok {
		return val, nil
	}
	return "", errors.New("empty")
}

func (c *CacheTest) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	c.mx.Lock()
	c.setCount++
	defer c.mx.Unlock()
	val, err := value.(model.Product).MarshalBinary()
	c.cache[key] = string(val)
	return err
}

func (c *CacheTest) Del(ctx context.Context, key string) error { return nil }

func TestProductCacheServiceWaitForCache(t *testing.T) {
	productServiceTest := &ProductServiceTest{}
	cache := &CacheTest{cache: make(map[string]string)}

	pcs := NewProductCacheService(productServiceTest, cache, 5*time.Second)

	ctx := context.Background()
	wg := &sync.WaitGroup{}

	testCount := 10000

	for range testCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pcs.GetProduct(ctx, 1)
		}()
	}

	wg.Wait()

	assert.Equal(t, 1, productServiceTest.counter, "productServiceTest.GetProduct")
	assert.Equal(t, 1, cache.setCount, "cacheTest.setCahce")
	assert.Equal(t, testCount*2-1, cache.getCount, "cacheTest.getCahce")
}
