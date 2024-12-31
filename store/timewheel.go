package store

import (
	"sync"
	"time"
)

// timeSlot 存储过期键的时间轮槽结构
type timeSlot struct {
	sync.Mutex
	keys map[string]struct{} // 存储键集合
}

// TimeWheel 是时间轮的核心结构
type TimeWheel struct {
	slots        []*timeSlot   // 槽数组
	slotCount    int           // 槽数
	tickInterval time.Duration // 每个槽位的时间间隔
	currentSlot  int           // 当前槽位置
	ticker       *time.Ticker  // 定时器
	stopChan     chan struct{} // 停止信号通道
	store        *MemoryStore  // 引用 MemoryStore
}

// NewTimeWheel 创建一个新的时间轮
func newTimeWheel(slotCount int, tickInterval time.Duration, store *MemoryStore) *TimeWheel {
	slots := make([]*timeSlot, slotCount)
	for i := 0; i < slotCount; i++ {
		slots[i] = &timeSlot{keys: make(map[string]struct{})}
	}

	return &TimeWheel{
		slots:        slots,
		slotCount:    slotCount,
		tickInterval: tickInterval,
		currentSlot:  0,
		ticker:       time.NewTicker(tickInterval),
		stopChan:     make(chan struct{}),
		store:        store,
	}
}

// Add 键加入到时间轮中
func (tw *TimeWheel) Add(key string, expireAt time.Time) {
	slotIndex := tw.calculateSlot(expireAt)

	slot := tw.slots[slotIndex]
	slot.Lock()
	defer slot.Unlock()

	// 将键添加到对应槽位
	slot.keys[key] = struct{}{}
}

// Remove 移除键
func (tw *TimeWheel) Remove(key string) {
	for _, slot := range tw.slots {
		slot.Lock()
		delete(slot.keys, key)
		slot.Unlock()
	}
}

// calculateSlot 计算键所属的槽位
func (tw *TimeWheel) calculateSlot(expireAt time.Time) int {
	duration := expireAt.Sub(time.Now())
	slotOffset := int(duration / tw.tickInterval)
	return (tw.currentSlot + slotOffset) % tw.slotCount
}

// Start 启动时间轮
func (tw *TimeWheel) Start() {
	for {
		select {
		case <-tw.ticker.C:
			tw.tick(tw.store)
		case <-tw.stopChan:
			tw.ticker.Stop()
			return
		}
	}
}

// tick 时间轮每个槽位更新时处理过期的键
func (tw *TimeWheel) tick(store *MemoryStore) {
	// 获取当前槽位
	currentSlot := tw.slots[tw.currentSlot]

	currentSlot.Lock()
	keys := make([]string, 0, len(currentSlot.keys))
	for key := range currentSlot.keys {
		keys = append(keys, key)
	}
	currentSlot.keys = make(map[string]struct{}) // 清空当前槽位
	currentSlot.Unlock()

	// 处理过期键
	for _, key := range keys {
		store.Delete(key) // 删除存储中的键
	}

	// 移动到下一个槽位
	tw.currentSlot = (tw.currentSlot + 1) % tw.slotCount
}

// Stop 停止时间轮
func (tw *TimeWheel) Stop() {
	close(tw.stopChan)
}
