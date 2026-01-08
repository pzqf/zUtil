package zCache

import (
	"testing"
	"time"
)

func TestLRUCache(t *testing.T) {
	cache := NewLRUCache(2, time.Second)
	defer cache.Stop()

	// 测试Set和Get
	err := cache.Set("key1", "value1", time.Second)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	value, err := cache.Get("key1")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected value1, got %v", value)
	}

	// 测试LRU淘汰
	err = cache.Set("key2", "value2", time.Second)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	err = cache.Set("key3", "value3", time.Second)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// key1应该被淘汰
	_, err = cache.Get("key1")
	if err == nil {
		t.Errorf("Expected key1 to be evicted")
	}

	// 测试Delete
	err = cache.Delete("key2")
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}
	_, err = cache.Get("key2")
	if err == nil {
		t.Errorf("Expected key2 to be deleted")
	}

	// 测试Clear
	err = cache.Clear()
	if err != nil {
		t.Errorf("Clear failed: %v", err)
	}
	if cache.Len() != 0 {
		t.Errorf("Expected cache to be empty, got %d items", cache.Len())
	}
}

func TestSimpleCache(t *testing.T) {
	cache := NewSimpleCache(time.Second)

	// 测试Set和Get
	err := cache.Set("key1", "value1", time.Second)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	value, err := cache.Get("key1")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected value1, got %v", value)
	}

	// 测试过期
	err = cache.Set("key2", "value2", 50*time.Millisecond)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	time.Sleep(100 * time.Millisecond)

	_, err = cache.Get("key2")
	if err == nil {
		t.Errorf("Expected key2 to be expired")
	}

	// 测试Delete
	err = cache.Delete("key1")
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}
	_, err = cache.Get("key1")
	if err == nil {
		t.Errorf("Expected key1 to be deleted")
	}

	// 测试Clear
	err = cache.Clear()
	if err != nil {
		t.Errorf("Clear failed: %v", err)
	}
	if cache.Len() != 0 {
		t.Errorf("Expected cache to be empty, got %d items", cache.Len())
	}
}

func TestTTLCache(t *testing.T) {
	cache := NewTTLCache(100 * time.Millisecond)
	defer cache.Stop()

	// 测试Set和Get
	err := cache.Set("key1", "value1", 0)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	value, err := cache.Get("key1")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected value1, got %v", value)
	}

	// 测试过期清理
	err = cache.Set("key2", "value2", 50*time.Millisecond)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
	time.Sleep(150 * time.Millisecond)

	// 等待清理完成
	time.Sleep(50 * time.Millisecond)

	// key2应该已被清理
	_, err = cache.Get("key2")
	if err == nil {
		t.Errorf("Expected key2 to be expired and cleaned up")
	}
}

func TestDefaultCache(t *testing.T) {
	// 测试默认缓存
	err := SetDefault("key1", "value1", time.Second)
	if err != nil {
		t.Errorf("SetDefault failed: %v", err)
	}

	value, err := GetDefault("key1")
	if err != nil {
		t.Errorf("GetDefault failed: %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected value1, got %v", value)
	}

	// 测试删除
	err = DeleteDefault("key1")
	if err != nil {
		t.Errorf("DeleteDefault failed: %v", err)
	}
	_, err = GetDefault("key1")
	if err == nil {
		t.Errorf("Expected key1 to be deleted")
	}
}

func TestCacheKeys(t *testing.T) {
	cache := NewLRUCache(3, time.Second)
	defer cache.Stop()

	cache.Set("key1", "value1", time.Second)
	cache.Set("key2", "value2", time.Second)
	cache.Set("key3", "value3", time.Second)

	keys := cache.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// 检查所有键都存在
	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}

	if !keyMap["key1"] || !keyMap["key2"] || !keyMap["key3"] {
		t.Errorf("Expected keys key1, key2, key3, got %v", keys)
	}
}

func TestCacheLen(t *testing.T) {
	cache := NewLRUCache(3, time.Second)
	defer cache.Stop()

	if cache.Len() != 0 {
		t.Errorf("Expected 0 items, got %d", cache.Len())
	}

	cache.Set("key1", "value1", time.Second)
	if cache.Len() != 1 {
		t.Errorf("Expected 1 item, got %d", cache.Len())
	}

	cache.Set("key2", "value2", time.Second)
	cache.Set("key3", "value3", time.Second)
	if cache.Len() != 3 {
		t.Errorf("Expected 3 items, got %d", cache.Len())
	}

	// 添加第4个，应该淘汰1个
	cache.Set("key4", "value4", time.Second)
	if cache.Len() != 3 {
		t.Errorf("Expected 3 items, got %d", cache.Len())
	}
}