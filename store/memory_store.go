package store

import (
	"hash/fnv"
	"sync"
	"time"
)

// shard 用于分片存储数据，每个分片维护独立锁和结构
type shard struct {
	sync.RWMutex
	data      map[string]any       // 数据存储
	expireMap map[string]time.Time // 存储每个键的过期时间
	numStored int                  // 当前存储数量
}

// MemoryStore 是主存储结构，包含多个分片
type MemoryStore struct {
	shards     []shard    // 数据分片
	shardCount int        // 分片数
	timeWheel  *TimeWheel // 时间轮实例
}

// NewMemoryStore 创建一个新的 MemoryStore
func NewMemoryStore(shardCount int, slotCount int, tickInterval time.Duration) *MemoryStore {
	shards := make([]shard, shardCount)
	for i := 0; i < shardCount; i++ {
		shards[i] = shard{
			data:      make(map[string]any),
			expireMap: make(map[string]time.Time),
		}
	}
	// 创建并启动时间轮
	timeWheel := newTimeWheel(slotCount, tickInterval, nil) // 这里传递 nil，稍后再设置 MemoryStore
	ms := &MemoryStore{
		shards:     shards,
		shardCount: shardCount,
		timeWheel:  timeWheel,
	}

	// 设置 timeWheel 的 store 引用
	ms.timeWheel.store = ms
	// 启动时间轮
	go ms.timeWheel.Start()

	return ms
}

// getShard 获取对应分片
func (ms *MemoryStore) getShard(key string) *shard {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	idx := int(hash.Sum32()) % ms.shardCount
	return &ms.shards[idx]
}

// Set 设置键值对，并处理过期时间
func (ms *MemoryStore) Set(key string, value any, ttl time.Duration) {
	shard := ms.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	// 设置值并记录过期时间
	if ttl == -1 {
		// 删除对应的过期key
		delete(shard.expireMap, key)
		ms.timeWheel.Remove(key)
	} else {
		expirationTime := time.Now().Add(ttl)
		shard.expireMap[key] = expirationTime // 记录过期时间
		// 通过时间轮添加键
		ms.timeWheel.Add(key, expirationTime)
	}

	shard.data[key] = value
	shard.numStored++
}

// Get 获取键值对，并检查是否过期
func (ms *MemoryStore) Get(key string, clear bool) (any, int64, bool) {
	shard := ms.getShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// 检查键是否存在
	value, exists := shard.data[key]
	if !exists {
		return nil, 0, false
	}

	// 检查过期时间
	nowTime := time.Now()
	expireAt, ok := shard.expireMap[key]
	if !ok || expireAt.Before(nowTime) {
		return nil, 0, false // 如果键过期，则返回 nil
	}

	// 清除指定 key
	if clear {
		ms.collectSpecifiedKey(shard, key)
	}

	seconds := int64(expireAt.Sub(nowTime).Seconds())
	return value, seconds, true
}

// collectSpecifiedKey 清除指定 key
func (ms *MemoryStore) collectSpecifiedKey(shard *shard, key string) {
	delete(shard.data, key)
	delete(shard.expireMap, key)
	shard.numStored--
}

// Delete 删除键
func (ms *MemoryStore) Delete(key string) {
	shard := ms.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	ms.collectSpecifiedKey(shard, key)
}

// IsExpired 检查指定键是否已过期
func (ms *MemoryStore) IsExpired(key string) bool {
	shard := ms.getShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// 检查键是否存在
	expireAt, exists := shard.expireMap[key]
	if !exists {
		return true // 键不存在，则视为已过期
	}

	// 检查是否过期
	return expireAt.Before(time.Now())
}

// Stats 返回当前统计信息
func (ms *MemoryStore) Stats() map[string]any {
	totalStored := 0
	for i := range ms.shards {
		totalStored += ms.shards[i].numStored
	}
	return map[string]any{
		"totalStored": totalStored,
		"shardCount":  ms.shardCount,
	}
}
