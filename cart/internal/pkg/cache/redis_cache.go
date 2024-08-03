package cache

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLRUCache struct {
	size         int
	cache        *redis.Client
	list         *list.List
	listElements map[string]*list.Element
	mx           sync.RWMutex
}

func NewRedisLRUCache(redisClient *redis.Client, size int) *RedisLRUCache {
	return &RedisLRUCache{
		size:         size,
		cache:        redisClient,
		list:         list.New(),
		listElements: make(map[string]*list.Element, size),
	}
}

func (c *RedisLRUCache) Get(ctx context.Context, key string) (string, error) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	val, err := c.cache.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("c.cache.Get: %w", err)
	}
	if el, ok := c.listElements[key]; ok {
		c.list.MoveToFront(el)
	}
	return val, nil
}

func (c *RedisLRUCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	if len(c.listElements) == c.size {
		lastKey := c.list.Back().Value.(string)
		if err := c.cache.Del(ctx, lastKey).Err(); err != nil {
			return fmt.Errorf("c.cache.Del: %w", err)
		}
		delete(c.listElements, lastKey)
		c.list.Remove(c.list.Back())
	}

	if err := c.cache.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("c.cache.Set: %w", err)
	}
	c.listElements[key] = c.list.PushFront(key)
	return nil
}

func (c *RedisLRUCache) Del(ctx context.Context, key string) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	if err := c.cache.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("c.cache.Del: %w", err)
	}
	if el, ok := c.listElements[key]; ok {
		c.list.Remove(el)
	}
	delete(c.listElements, key)
	return nil
}
