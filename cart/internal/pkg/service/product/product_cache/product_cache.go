package product_cache

import (
	"context"
	"fmt"
	"route256/cart/internal/pkg/model"
	"route256/cart/internal/pkg/utils/metrics"
	"route256/cart/pkg/tracing"
	"strconv"
	"sync"
	"time"
)

type ProductService interface {
	GetProduct(context.Context, model.ProductSku) (*model.Product, error)
}

type Cacher interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

const (
	getProductCacheKey = "cart:product:get_product"
)

type ProductCacheService struct {
	productService ProductService
	cache          Cacher
	defaultTTL     time.Duration
	inProcess      map[string]chan struct{}
	inProcessMx    sync.RWMutex
}

func NewProductCacheService(productService ProductService, cache Cacher, defaultTTL time.Duration) *ProductCacheService {
	return &ProductCacheService{
		productService: productService,
		cache:          cache,
		defaultTTL:     defaultTTL,
		inProcess:      make(map[string]chan struct{}),
	}
}

func (p *ProductCacheService) GetProduct(ctx context.Context, ProductSku model.ProductSku) (_ *model.Product, err error) {
	ctx, span := tracing.Start(ctx, "ProductCache.GetProduct")
	defer tracing.EndWithCheckError(span, &err)

	serviceHandler := "product.GetProduct"

	start := time.Now()

	key := getProductCacheKey + ":" + strconv.Itoa(int(ProductSku))

	if cacheProduct, err := p.cache.Get(ctx, key); err == nil {
		go func(start time.Time) {
			metrics.CacheHitCounter(serviceHandler)
			metrics.CacheHitDuration(serviceHandler, time.Since(start).Seconds())
		}(start)

		var product model.Product
		if err := product.UnmarshalBinary([]byte(cacheProduct)); err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w", err)
		}
		return &product, nil
	}

	p.inProcessMx.RLock()
	if done, ok := p.inProcess[key]; ok {
		p.inProcessMx.RUnlock()
		<-done
		return p.GetProduct(ctx, ProductSku)
	}
	p.inProcessMx.RUnlock()
	p.inProcessMx.Lock()
	p.inProcess[key] = make(chan struct{})
	p.inProcessMx.Unlock()
	defer func() {
		p.inProcessMx.Lock()
		close(p.inProcess[key])
		delete(p.inProcess, key)
		p.inProcessMx.Unlock()
	}()

	product, err := p.productService.GetProduct(ctx, ProductSku)
	if err != nil {
		return nil, fmt.Errorf("p.productService.GetProduct: %w", err)
	}

	go func(start time.Time) {
		metrics.CacheMissCounter(serviceHandler)
		metrics.CacheMissDuration(serviceHandler, time.Since(start).Seconds())
	}(start)

	err = p.cache.Set(ctx, key, *product, p.defaultTTL)
	if err != nil {
		return nil, fmt.Errorf("p.cache.Set: %w", err)
	}
	return product, nil
}
