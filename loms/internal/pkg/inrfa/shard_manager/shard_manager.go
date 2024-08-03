package shard_manager

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spaolacci/murmur3"
)

var (
	ErrShardIndexOutOfRange = errors.New("shard index is out of range")
)

type ShardKey string
type ShardIndex int
type ShardFn func(ShardKey) ShardIndex

func GetMurmur3ShardFn(shardsCnt int) ShardFn {
	hasher := murmur3.New32()
	return func(key ShardKey) ShardIndex {
		defer hasher.Reset()
		_, _ = hasher.Write([]byte(key))
		return ShardIndex(hasher.Sum32() % uint32(shardsCnt))
	}
}

type ShardManager struct {
	fn     ShardFn
	shards []*pgxpool.Pool
}

func NewShardManager(fn ShardFn, shards ...*pgxpool.Pool) *ShardManager {
	return &ShardManager{
		fn:     fn,
		shards: shards,
	}
}

func (m *ShardManager) GetShardIndex(key ShardKey) ShardIndex {
	return m.fn(key)
}

func (m *ShardManager) GetShardIndexFromID(id int64) ShardIndex {
	return ShardIndex(id % 1000)
}

func (m *ShardManager) Pick(index ShardIndex) (*pgxpool.Pool, error) {
	if int(index) < len(m.shards) {
		return m.shards[index], nil
	}
	return nil, fmt.Errorf("%w: given index=%d, len=%d", ErrShardIndexOutOfRange, index, len(m.shards))
}

func (m *ShardManager) GetShards() []*pgxpool.Pool {
	return m.shards
}
