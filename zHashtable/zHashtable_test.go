package zHashtable

import (
	"fmt"
	"testing"
)

// TestNewHashTable 测试创建哈希表
func TestNewHashTable(t *testing.T) {
	ht := NewHashTable()
	if ht == nil {
		t.Fatal("NewHashTable() returned nil")
	}
	if ht.capacity != 256 {
		t.Errorf("Expected capacity 256, got %d", ht.capacity)
	}
	if ht.size != 0 {
		t.Errorf("Expected size 0, got %d", ht.size)
	}
}

// TestAddAndGet 测试添加和获取元素
func TestAddAndGet(t *testing.T) {
	ht := NewHashTable()

	// 添加元素
	ht.Add("key1", "value1")
	ht.Add("key2", 123)
	ht.Add("key3", true)

	// 获取元素
	val1, exists1 := ht.Get("key1")
	if !exists1 {
		t.Error("Expected key1 to exist")
	}
	if val1 != "value1" {
		t.Errorf("Expected value1 for key1, got %v", val1)
	}

	val2, exists2 := ht.Get("key2")
	if !exists2 {
		t.Error("Expected key2 to exist")
	}
	if val2 != 123 {
		t.Errorf("Expected 123 for key2, got %v", val2)
	}

	val3, exists3 := ht.Get("key3")
	if !exists3 {
		t.Error("Expected key3 to exist")
	}
	if val3 != true {
		t.Errorf("Expected true for key3, got %v", val3)
	}

	// 获取不存在的键
	val4, exists4 := ht.Get("key4")
	if exists4 {
		t.Error("Expected key4 to not exist")
	}
	if val4 != nil {
		t.Errorf("Expected nil for key4, got %v", val4)
	}
}

// TestSet 测试更新元素
func TestSet(t *testing.T) {
	ht := NewHashTable()

	// 先添加元素
	ht.Add("key1", "value1")

	// 更新已存在的键
	updated := ht.Set("key1", "new_value1")
	if !updated {
		t.Error("Expected Set to return true for existing key")
	}

	// 检查更新后的值
	val, exists := ht.Get("key1")
	if !exists {
		t.Error("Expected key1 to exist after Set")
	}
	if val != "new_value1" {
		t.Errorf("Expected new_value1 for key1 after Set, got %v", val)
	}

	// 更新不存在的键
	updated2 := ht.Set("key2", "value2")
	if updated2 {
		t.Error("Expected Set to return false for non-existing key")
	}

	// 检查不存在的键是否真的不存在
	val2, exists2 := ht.Get("key2")
	if exists2 {
		t.Error("Expected key2 to not exist after failed Set")
	}
	if val2 != nil {
		t.Errorf("Expected nil for key2 after failed Set, got %v", val2)
	}
}

// TestRemove 测试删除元素
func TestRemove(t *testing.T) {
	ht := NewHashTable()

	// 添加元素
	ht.Add("key1", "value1")
	ht.Add("key2", "value2")

	// 删除存在的键
	removed := ht.Remove("key1")
	if !removed {
		t.Error("Expected Remove to return true for existing key")
	}

	// 检查是否真的删除了
	val, exists := ht.Get("key1")
	if exists {
		t.Error("Expected key1 to not exist after Remove")
	}
	if val != nil {
		t.Errorf("Expected nil for key1 after Remove, got %v", val)
	}

	// 检查其他元素是否还存在
	val2, exists2 := ht.Get("key2")
	if !exists2 {
		t.Error("Expected key2 to still exist after removing key1")
	}
	if val2 != "value2" {
		t.Errorf("Expected value2 for key2 after removing key1, got %v", val2)
	}

	// 删除不存在的键
	removed2 := ht.Remove("key3")
	if removed2 {
		t.Error("Expected Remove to return false for non-existing key")
	}
}

// TestResize 测试哈希表自动扩容
func TestResize(t *testing.T) {
	ht := NewHashTable()

	// 获取初始容量
	initialCapacity := ht.capacity

	// 添加足够多的元素以触发扩容
	for i := 0; i < initialCapacity; i++ {
		ht.Add(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	// 检查是否扩容
	if ht.capacity <= initialCapacity {
		t.Errorf("Expected capacity to increase from %d, got %d", initialCapacity, ht.capacity)
	}

	// 验证所有元素仍然可以访问
	for i := 0; i < initialCapacity; i++ {
		key := fmt.Sprintf("key%d", i)
		val, exists := ht.Get(key)
		if !exists {
			t.Errorf("Expected key %s to exist after resize", key)
		}
		if val != fmt.Sprintf("value%d", i) {
			t.Errorf("Expected value %s for key %s after resize, got %v", fmt.Sprintf("value%d", i), key, val)
		}
	}
}
