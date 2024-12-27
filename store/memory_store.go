package store

import (
	"container/list"
	"sync"
	"time"
)

type idByTimeValue struct {
	timestamp time.Time
	id        string
}

type memoryStore struct {
	sync.RWMutex
	data       map[string]any           // 容器
	expireList *list.List               // 双向链表，按时间排序存储 id
	expireMap  map[string]*list.Element //
	expiration time.Duration            // 过期时间
	threshold  int                      // 清理阈值
	numStored  int                      // 当前存储的键值对数
}

func NewMemoryStore(threshold int, expiration time.Duration) MemoryStore {
	return &memoryStore{
		data:       make(map[string]any),
		expireList: list.New(),
		expireMap:  make(map[string]*list.Element),
		expiration: expiration,
		threshold:  threshold,
		numStored:  0,
	}
}

func (ms *memoryStore) Set(key string, value any) {
	ms.Lock()
	ms.data[key] = value
	ms.expireList.PushBack(idByTimeValue{timestamp: time.Now(), id: key})
	ms.numStored++
	needCollect := ms.numStored > ms.collectNum
	ms.Unlock()
	if needCollect {
		go ms.collect() // 使用 goroutine 异步清理
	}
}

func (ms *memoryStore) Get(key string, clear bool) (any, bool) {
	ms.RLock()
	value, exists := ms.data[key]
	ms.RUnlock()
	if !exists {
		return nil, false
	}

	if clear {
		ms.Lock()
		defer ms.Unlock()
		e := ms.searchEle(key)
		ms.collectSpecifiedKey(e)
	}
	return value, true
}

func (ms *memoryStore) Delete(key string) {
	ms.Lock()
	defer ms.Unlock()
	delete(ms.data, key)
}

// Cleanup 清空容器里的key, 默认清除所有, 通过onlyExpired可选是否清除仅过期的
func (ms *memoryStore) Cleanup(onlyExpired bool) {
	if onlyExpired {
		ms.collect()
	} else {
		ms.Lock()
		defer ms.Unlock()
		ms.data = make(map[string]any) // 直接重新初始化map
		ms.expireList.Init()           // 清空双向链表
		ms.numStored = 0               // 重置存储计数
	}
}

func (ms *memoryStore) searchEle(key string) *list.Element {
	for e := ms.expireList.Front(); e != nil; e = e.Next() {
		if ev, ok := e.Value.(idByTimeValue); ok && ev.id == key {
			return e
		}
	}
}

func (ms *memoryStore) collect() {
	ms.Lock()
	defer ms.Unlock()

	now := time.Now()
	for e := ms.expireList.Front(); e != nil; e = e.Next() {
		ms.collectSpecifiedTime(e, now)
	}
}

// collect 清除指定过期时间的key
func (ms *memoryStore) collectSpecifiedTime(e *list.Element, specifyTime time.Time) {
	ev, ok := e.Value.(idByTimeValue)
	if ok && ev.timestamp.Add(ms.expiration).Before(specifyTime) {
		delete(ms.data, ev.id) // 删除过期键
		ms.expireList.Remove(e)
		ms.numStored--
	}
}

// collectSpecifiedKey 清除指定的KEY
func (ms *memoryStore) collectSpecifiedKey(e *list.Element) {
	ev, ok := e.Value.(idByTimeValue)
	if ok {
		delete(ms.data, ev.id) // 删除过期键
		ms.expireList.Remove(e)
		ms.numStored--
	}
}
