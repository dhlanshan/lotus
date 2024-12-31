package store

import (
	"fmt"
	"testing"
	"time"
)

// 测试 MemoryStore 的 Set 和 Get 方法
func TestMemoryStore_Set_Get(t *testing.T) {
	ms := NewMemoryStore(4, 10, time.Second) // 创建一个带有 4 个分片的 MemoryStore

	// 设置一个有效的键值对，TTL 设置为 2 秒
	ms.Set("key1", "value1", 2*time.Second)

	// 尝试获取该键
	value, ttl, exists := ms.Get("key1", false)
	if !exists {
		t.Errorf("Expected key1 to be found, but got not found")
	}
	if value != "value1" {
		t.Errorf("Expected value1, but got %v", value)
	}
	fmt.Println("---", value)
	fmt.Println("---ttl", ttl)

	// 等待键过期
	time.Sleep(3 * time.Second)

	// 尝试获取过期的键
	value, ttl, exists = ms.Get("key1", false)
	if exists {
		t.Errorf("Expected key1 to be expired, but it is still found")
	}
	fmt.Println("---", value)
	fmt.Println("---", ttl)
}

// 测试 MemoryStore 的 Delete 方法
func TestMemoryStore_Delete(t *testing.T) {
	ms := NewMemoryStore(4, 10, time.Second)

	// 设置键值对
	ms.Set("key1", "value1", 2*time.Second)

	// 删除键
	ms.Delete("key1")

	// 尝试获取已删除的键
	_, _, exists := ms.Get("key1", false)
	if exists {
		t.Errorf("Expected key1 to be deleted, but it is still found")
	}
}

// 测试 MemoryStore 的 IsExpired 方法
func TestMemoryStore_IsExpired(t *testing.T) {
	ms := NewMemoryStore(4, 10, time.Second)

	// 设置一个有效的键值对，TTL 设置为 2 秒
	ms.Set("key1", "value1", 2*time.Second)

	// 等待一秒，键应该未过期
	time.Sleep(1 * time.Second)
	if ms.IsExpired("key1") {
		t.Errorf("Expected key1 to not be expired, but it is expired")
	}

	// 等待超过 TTL，键应该过期
	time.Sleep(2 * time.Second)
	if !ms.IsExpired("key1") {
		t.Errorf("Expected key1 to be expired, but it is not expired")
	}
}

// 测试 MemoryStore 的 Stats 方法
func TestMemoryStore_Stats(t *testing.T) {
	ms := NewMemoryStore(4, 10, time.Second)

	// 设置两个键值对
	ms.Set("key1", "value1", 2*time.Second)
	ms.Set("key2", "value2", 2*time.Second)

	// 获取存储的统计信息
	stats := ms.Stats()
	if stats["totalStored"] != 2 {
		t.Errorf("Expected totalStored to be 2, but got %v", stats["totalStored"])
	}
}
