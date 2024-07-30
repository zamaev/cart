package cache

import (
	"context"
	"route256/cart/internal/pkg/config"
	"route256/cart/internal/pkg/model"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getCache(t *testing.T) *RedisLRUCache {
	config := config.NewConfig()
	redisCache := redis.NewClient(&redis.Options{
		Addr:     config.RedisUrl,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})
	_, err := redisCache.Ping(context.Background()).Result()
	require.NoError(t, err)
	return NewRedisLRUCache(redisCache, 3)
}

func TestCache(t *testing.T) {
	ctx := context.Background()
	cache := getCache(t)

	product := model.Product{
		Sku:  1,
		Name: "test",
	}

	err := cache.Set(ctx, "1", product, 0)
	assert.NoError(t, err)

	res, err := cache.Get(ctx, "1")
	assert.NoError(t, err)

	var cahceProduct model.Product
	err = cahceProduct.UnmarshalBinary([]byte(res))
	assert.NoError(t, err)

	assert.Equal(t, cahceProduct, product)
}

func TestCacheSize(t *testing.T) {
	ctx := context.Background()
	cache := getCache(t)

	err := cache.Set(ctx, "1", "toy1", 0)
	assert.NoError(t, err)

	toy1, err := cache.Get(ctx, "1")
	assert.NoError(t, err)
	assert.Equal(t, "toy1", toy1)

	cache.Set(ctx, "2", "toy2", 0)
	cache.Set(ctx, "3", "toy3", 0)
	cache.Set(ctx, "4", "toy4", 0)

	// после добавления 4, вытесняется 1
	_, err = cache.Get(ctx, "1")
	assert.Error(t, err)

	toy4, err := cache.Get(ctx, "4")
	assert.NoError(t, err)
	assert.Equal(t, "toy4", toy4)

	err = cache.Del(ctx, "4")
	assert.NoError(t, err)

	_, err = cache.Get(ctx, "4")
	assert.Error(t, err)
}
