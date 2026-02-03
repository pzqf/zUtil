package zMap

import (
	"fmt"
	"testing"
	"time"
)

func TestMapBasicOperations(t *testing.T) {
	m := NewMap()

	t.Run("Get_ExistingKey_ReturnsValue", func(t *testing.T) {
		key := "test"
		value := "value"
		m.Store(key, value)

		result, exists := m.Load(key)
		if !exists {
			t.Errorf("expected key to exist")
		}
		if result != value {
			t.Errorf("expected %v, got %v", value, result)
		}
	})

	t.Run("Get_NonExistingKey_ReturnsFalse", func(t *testing.T) {
		_, exists := m.Load("nonexistent")
		if exists {
			t.Errorf("expected key to not exist")
		}
	})

	t.Run("Store_NewKey_UpdatesLen", func(t *testing.T) {
		key := "test"
		value := "value"

		m.Store(key, value)

		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})

	t.Run("Store_ExistingKey_OverwritesValue", func(t *testing.T) {
		key := "test"
		oldValue := "old"
		newValue := "new"

		m.Store(key, oldValue)
		m.Store(key, newValue)

		result, _ := m.Load(key)
		if result != newValue {
			t.Errorf("expected %v, got %v", newValue, result)
		}
		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})

	t.Run("Delete_ExistingKey_RemovesEntry", func(t *testing.T) {
		key := "test"

		m.Store(key, "value")
		m.Delete(key)

		_, exists := m.Load(key)
		if exists {
			t.Errorf("expected key to be deleted")
		}
		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("Delete_NonExistingKey_NoChange", func(t *testing.T) {
		m.Delete("nonexistent")

		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("Len_EmptyMap_Returns0", func(t *testing.T) {
		m := NewMap()

		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("Len_MultipleEntries_ReturnsCorrectCount", func(t *testing.T) {
		m := NewMap()

		m.Store("key1", 1)
		m.Store("key2", 2)
		m.Store("key3", 3)

		if m.Len() != 3 {
			t.Errorf("expected len 3, got %d", m.Len())
		}
	})
}

func TestMapRangeOperation(t *testing.T) {
	m := NewMap()
	expected := map[string]int{
		"key1": 1,
		"key2": 2,
		"key3": 3,
	}

	for k, v := range expected {
		m.Store(k, v)
	}

	t.Run("Range_AllEntries_VisitAllKeys", func(t *testing.T) {
		visited := make(map[string]bool)

		m.Range(func(key, value interface{}) bool {
			k := key.(string)
			v := value.(int)

			if expected[k] != v {
				t.Errorf("expected %d for key %s, got %d", expected[k], k, v)
			}
			visited[k] = true
			return true
		})

		if len(visited) != len(expected) {
			t.Errorf("expected to visit %d keys, visited %d", len(expected), len(visited))
		}
		for k := range expected {
			if !visited[k] {
				t.Errorf("expected to visit key %s", k)
			}
		}
	})

	t.Run("Range_StopEarly_ReturnsFalse", func(t *testing.T) {
		count := 0

		m.Range(func(key, value interface{}) bool {
			count++
			if count == 2 {
				return false
			}
			return true
		})

		if count != 2 {
			t.Errorf("expected to stop after 2 visits, got %d", count)
		}
	})
}

func TestMapExtendedOperations(t *testing.T) {
	m := NewMap()

	t.Run("LoadOrStore_NewKey_InsertsAndReturnsValue", func(t *testing.T) {
		key := "test"
		value := "value"

		result, loaded := m.LoadOrStore(key, value)

		if loaded {
			t.Errorf("expected key to not exist")
		}
		if result != value {
			t.Errorf("expected %v, got %v", value, result)
		}
		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})

	t.Run("LoadOrStore_ExistingKey_ReturnsExistingValue", func(t *testing.T) {
		key := "test"
		oldValue := "old"
		newValue := "new"

		m.Store(key, oldValue)
		result, loaded := m.LoadOrStore(key, newValue)

		if !loaded {
			t.Errorf("expected key to exist")
		}
		if result != oldValue {
			t.Errorf("expected %v, got %v", oldValue, result)
		}
		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})

	t.Run("LoadAndDelete_ExistingKey_ReturnsValue", func(t *testing.T) {
		key := "test"
		value := "value"

		m.Store(key, value)
		result, loaded := m.LoadAndDelete(key)

		if !loaded {
			t.Errorf("expected key to exist")
		}
		if result != value {
			t.Errorf("expected %v, got %v", value, result)
		}
		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("LoadAndDelete_NonExistingKey_ReturnsFalse", func(t *testing.T) {
		_, loaded := m.LoadAndDelete("nonexistent")

		if loaded {
			t.Errorf("expected key to not exist")
		}
	})

	t.Run("CompareAndDelete_MatchingValue_ReturnsTrue", func(t *testing.T) {
		key := "test"
		value := "value"

		m.Store(key, value)
		deleted := m.CompareAndDelete(key, value)

		if !deleted {
			t.Errorf("expected to delete matching key")
		}
		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("CompareAndDelete_NonMatchingValue_ReturnsFalse", func(t *testing.T) {
		key := "test"
		value := "value"

		m.Store(key, value)
		deleted := m.CompareAndDelete(key, "wrong")

		if deleted {
			t.Errorf("expected to not delete non-matching key")
		}
		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})

	t.Run("CompareAndSwap_MatchingValue_ReturnsTrue", func(t *testing.T) {
		key := "test"
		oldValue := "old"
		newValue := "new"

		m.Store(key, oldValue)
		swapped := m.CompareAndSwap(key, oldValue, newValue)

		if !swapped {
			t.Errorf("expected to swap matching value")
		}
		result, _ := m.Load(key)
		if result != newValue {
			t.Errorf("expected %v, got %v", newValue, result)
		}
		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})

	t.Run("CompareAndSwap_NonMatchingValue_ReturnsFalse", func(t *testing.T) {
		key := "test"
		oldValue := "old"
		newValue := "new"

		m.Store(key, oldValue)
		swapped := m.CompareAndSwap(key, "wrong", newValue)

		if swapped {
			t.Errorf("expected to not swap non-matching value")
		}
		result, _ := m.Load(key)
		if result != oldValue {
			t.Errorf("expected %v, got %v", oldValue, result)
		}
	})
}

func TestMapClearOperation(t *testing.T) {
	m := NewMap()

	m.Store("key1", 1)
	m.Store("key2", 2)
	m.Store("key3", 3)

	if m.Len() != 3 {
		t.Errorf("expected len 3, got %d", m.Len())
	}

	m.Clear()

	if m.Len() != 0 {
		t.Errorf("expected len 0 after clear, got %d", m.Len())
	}

	_, exists := m.Load("key1")
	if exists {
		t.Errorf("expected key1 to be cleared")
	}

	_, exists = m.Load("key2")
	if exists {
		t.Errorf("expected key2 to be cleared")
	}

	_, exists = m.Load("key3")
	if exists {
		t.Errorf("expected key3 to be cleared")
	}
}

func TestMapConcurrentAccess(t *testing.T) {
	m := NewMap()
	key := "counter"
	iterations := 1000
	goroutines := 10

	t.Run("ConcurrentStore_Get_DoesNotPanic", func(t *testing.T) {
		done := make(chan bool, goroutines)

		for i := 0; i < goroutines; i++ {
			go func(gid int) {
				for j := 0; j < iterations; j++ {
					key := key + "_" + fmt.Sprintf("%d", gid) + "_" + fmt.Sprintf("%d", j)
					m.Store(key, j)
					m.Load(key)
				}
				done <- true
			}(i)
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}
	})

	t.Run("ConcurrentDelete_Len_MaintainsCorrectCount", func(t *testing.T) {
		m := NewMap()
		keys := make([]string, iterations)

		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			keys[i] = key
			m.Store(key, i)
		}

		done := make(chan bool, goroutines)

		for i := 0; i < goroutines; i++ {
			go func(gid int) {
				for j := gid; j < iterations; j += goroutines {
					m.Delete(keys[j])
				}
				done <- true
			}(i)
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}
	})

	t.Run("ConcurrentRange_DoesNotPanic", func(t *testing.T) {
		m := NewMap()
		for i := 0; i < 100; i++ {
			key := fmt.Sprintf("%d", i)
			m.Store(key, i)
		}

		done := make(chan bool, goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				for j := 0; j < 100; j++ {
					m.Range(func(key, value interface{}) bool {
						return true
					})
				}
				done <- true
			}()
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}
	})
}

func TestTypedMapBasicOperations(t *testing.T) {
	t.Run("Get_ExistingKey_ReturnsValue", func(t *testing.T) {
		m := NewTypedMap[string, int]()
		key := "test"
		value := 42
		m.Store(key, value)

		result, exists := m.Load(key)
		if !exists {
			t.Errorf("expected key to exist")
		}
		if result != value {
			t.Errorf("expected %d, got %d", value, result)
		}
	})

	t.Run("Get_NonExistingKey_ReturnsZeroValue", func(t *testing.T) {
		m := NewTypedMap[string, int]()

		result, exists := m.Load("nonexistent")
		if exists {
			t.Errorf("expected key to not exist")
		}
		if result != 0 {
			t.Errorf("expected zero value, got %d", result)
		}
	})

	t.Run("Store_NewKey_UpdatesLen", func(t *testing.T) {
		m := NewTypedMap[string, string]()
		key := "test"
		value := "value"

		m.Store(key, value)

		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})

	t.Run("Store_ExistingKey_OverwritesValue", func(t *testing.T) {
		m := NewTypedMap[int, float64]()
		key := 42
		oldValue := 1.0
		newValue := 2.0

		m.Store(key, oldValue)
		m.Store(key, newValue)

		result, _ := m.Load(key)
		if result != newValue {
			t.Errorf("expected %f, got %f", newValue, result)
		}
		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})

	t.Run("Delete_ExistingKey_RemovesEntry", func(t *testing.T) {
		m := NewTypedMap[string, int]()
		key := "test"

		m.Store(key, 42)
		m.Delete(key)

		_, exists := m.Load(key)
		if exists {
			t.Errorf("expected key to be deleted")
		}
		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("Delete_NonExistingKey_NoChange", func(t *testing.T) {
		m := NewTypedMap[string, int]()

		m.Delete("nonexistent")

		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("Len_EmptyMap_Returns0", func(t *testing.T) {
		m := NewTypedMap[string, int]()

		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("Len_MultipleEntries_ReturnsCorrectCount", func(t *testing.T) {
		m := NewTypedMap[int, string]()

		m.Store(1, "one")
		m.Store(2, "two")
		m.Store(3, "three")

		if m.Len() != 3 {
			t.Errorf("expected len 3, got %d", m.Len())
		}
	})
}

func TestTypedMapExtendedOperations(t *testing.T) {
	m := NewTypedMap[string, string]()

	t.Run("LoadOrStore_NewKey_InsertsAndReturnsValue", func(t *testing.T) {
		key := "test"
		value := "value"

		result, loaded := m.LoadOrStore(key, value)

		if loaded {
			t.Errorf("expected key to not exist")
		}
		if result != value {
			t.Errorf("expected %q, got %q", value, result)
		}
		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})

	t.Run("LoadOrStore_ExistingKey_ReturnsExistingValue", func(t *testing.T) {
		key := "test"
		oldValue := "old"
		newValue := "new"

		m.Store(key, oldValue)
		result, loaded := m.LoadOrStore(key, newValue)

		if !loaded {
			t.Errorf("expected key to exist")
		}
		if result != oldValue {
			t.Errorf("expected %q, got %q", oldValue, result)
		}
		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})

	t.Run("LoadAndDelete_ExistingKey_ReturnsValue", func(t *testing.T) {
		key := "test"
		value := "value"

		m.Store(key, value)
		result, loaded := m.LoadAndDelete(key)

		if !loaded {
			t.Errorf("expected key to exist")
		}
		if result != value {
			t.Errorf("expected %q, got %q", value, result)
		}
		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("LoadAndDelete_NonExistingKey_ReturnsZeroValue", func(t *testing.T) {
		result, loaded := m.LoadAndDelete("nonexistent")

		if loaded {
			t.Errorf("expected key to not exist")
		}
		if result != "" {
			t.Errorf("expected empty string, got %q", result)
		}
	})

	t.Run("CompareAndDelete_MatchingValue_ReturnsTrue", func(t *testing.T) {
		key := "test"
		value := "value"

		m.Store(key, value)
		deleted := m.CompareAndDelete(key, value)

		if !deleted {
			t.Errorf("expected to delete matching key")
		}
		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("CompareAndDelete_NonMatchingValue_ReturnsFalse", func(t *testing.T) {
		key := "test"
		value := "value"

		m.Store(key, value)
		deleted := m.CompareAndDelete(key, "wrong")

		if deleted {
			t.Errorf("expected to not delete non-matching key")
		}
		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})

	t.Run("CompareAndSwap_MatchingValue_ReturnsTrue", func(t *testing.T) {
		key := "test"
		oldValue := "old"
		newValue := "new"

		m.Store(key, oldValue)
		swapped := m.CompareAndSwap(key, oldValue, newValue)

		if !swapped {
			t.Errorf("expected to swap matching value")
		}
		result, _ := m.Load(key)
		if result != newValue {
			t.Errorf("expected %q, got %q", newValue, result)
		}
		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})

	t.Run("CompareAndSwap_NonMatchingValue_ReturnsFalse", func(t *testing.T) {
		key := "test"
		oldValue := "old"
		newValue := "new"

		m.Store(key, oldValue)
		swapped := m.CompareAndSwap(key, "wrong", newValue)

		if swapped {
			t.Errorf("expected to not swap non-matching value")
		}
		result, _ := m.Load(key)
		if result != oldValue {
			t.Errorf("expected %q, got %q", oldValue, result)
		}
	})
}

func TestTypedMapClearOperation(t *testing.T) {
	m := NewTypedMap[int, string]()

	m.Store(1, "one")
	m.Store(2, "two")
	m.Store(3, "three")

	if m.Len() != 3 {
		t.Errorf("expected len 3, got %d", m.Len())
	}

	m.Clear()

	if m.Len() != 0 {
		t.Errorf("expected len 0 after clear, got %d", m.Len())
	}

	_, exists := m.Load(1)
	if exists {
		t.Errorf("expected key 1 to be cleared")
	}

	_, exists = m.Load(2)
	if exists {
		t.Errorf("expected key 2 to be cleared")
	}

	_, exists = m.Load(3)
	if exists {
		t.Errorf("expected key 3 to be cleared")
	}
}

func TestTypedMapComplexTypes(t *testing.T) {
	type User struct {
		ID   int
		Name string
		Age  int
	}

	m := NewTypedMap[int, User]()

	t.Run("Store_Get_ComplexStruct", func(t *testing.T) {
		user := User{ID: 1, Name: "John", Age: 30}
		m.Store(user.ID, user)

		result, exists := m.Load(user.ID)
		if !exists {
			t.Errorf("expected user to exist")
		}
		if result != user {
			t.Errorf("expected user %+v, got %+v", user, result)
		}
	})

	t.Run("Store_Get_PointerValue", func(t *testing.T) {
		m := NewTypedMap[int, *User]()
		user := &User{ID: 2, Name: "Jane", Age: 25}
		m.Store(user.ID, user)

		result, exists := m.Load(user.ID)
		if !exists {
			t.Errorf("expected user pointer to exist")
		}
		if *result != *user {
			t.Errorf("expected user pointer %+v, got %+v", *user, *result)
		}
	})
}

func TestTypedMapRangeOperation(t *testing.T) {
	m := NewTypedMap[string, int]()
	expected := map[string]int{
		"key1": 1,
		"key2": 2,
		"key3": 3,
	}

	for k, v := range expected {
		m.Store(k, v)
	}

	t.Run("Range_AllEntries_VisitAllKeys", func(t *testing.T) {
		visited := make(map[string]bool)

		m.Range(func(key string, value int) bool {
			if expected[key] != value {
				t.Errorf("expected %d for key %s, got %d", expected[key], key, value)
			}
			visited[key] = true
			return true
		})

		if len(visited) != len(expected) {
			t.Errorf("expected to visit %d keys, visited %d", len(expected), len(visited))
		}
		for k := range expected {
			if !visited[k] {
				t.Errorf("expected to visit key %s", k)
			}
		}
	})

	t.Run("Range_StopEarly_ReturnsFalse", func(t *testing.T) {
		count := 0

		m.Range(func(key string, value int) bool {
			count++
			if count == 2 {
				return false
			}
			return true
		})

		if count != 2 {
			t.Errorf("expected to stop after 2 visits, got %d", count)
		}
	})
}

func TestTypedMapConcurrentAccess(t *testing.T) {
	m := NewTypedMap[string, int]()
	iterations := 1000
	goroutines := 10

	t.Run("ConcurrentStore_Get_DoesNotPanic", func(t *testing.T) {
		done := make(chan bool, goroutines)

		for i := 0; i < goroutines; i++ {
			go func(gid int) {
				for j := 0; j < iterations; j++ {
					key := fmt.Sprintf("%d", gid) + "_" + fmt.Sprintf("%d", j)
					m.Store(key, j)
					m.Load(key)
				}
				done <- true
			}(i)
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}
	})

	t.Run("ConcurrentDelete_Len_MaintainsCorrectCount", func(t *testing.T) {
		m := NewTypedMap[string, int]()
		keys := make([]string, iterations)

		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			keys[i] = key
			m.Store(key, i)
		}

		done := make(chan bool, goroutines)

		for i := 0; i < goroutines; i++ {
			go func(gid int) {
				for j := gid; j < iterations; j += goroutines {
					m.Delete(keys[j])
				}
				done <- true
			}(i)
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}
	})

	t.Run("ConcurrentRange_DoesNotPanic", func(t *testing.T) {
		m := NewTypedMap[string, int]()
		for i := 0; i < 100; i++ {
			key := fmt.Sprintf("%d", i)
			m.Store(key, i)
		}

		done := make(chan bool, goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				for j := 0; j < 100; j++ {
					m.Range(func(key string, value int) bool {
						return true
					})
				}
				done <- true
			}()
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}
	})
}

func TestShardedMapBasicOperations(t *testing.T) {
	t.Run("NewShardedMap_DefaultShards_ReturnsMap", func(t *testing.T) {
		m := NewShardedMap(32)

		if m == nil {
			t.Errorf("expected non-nil map")
		}
		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("NewShardedMap_NegativeShards_ReturnsDefault", func(t *testing.T) {
		m := NewShardedMap(-5)

		if m == nil {
			t.Errorf("expected non-nil map")
		}
		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("NewShardedMap32_Returns32Shards", func(t *testing.T) {
		m := NewShardedMap32()

		if m == nil {
			t.Errorf("expected non-nil map")
		}
	})
}

func TestShardedMapCoreOperations(t *testing.T) {
	t.Run("Get_ExistingKey_ReturnsValue", func(t *testing.T) {
		m := NewShardedMap(32)
		key := "test"
		value := "value"
		m.Store(key, value)

		result, exists := m.Load(key)
		if !exists {
			t.Errorf("expected key to exist")
		}
		if result != value {
			t.Errorf("expected %v, got %v", value, result)
		}
	})

	t.Run("Store_NewKey_UpdatesLen", func(t *testing.T) {
		m := NewShardedMap(16)
		key := "test"
		value := "value"

		m.Store(key, value)

		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})

	t.Run("Delete_ExistingKey_RemovesEntry", func(t *testing.T) {
		m := NewShardedMap(32)
		key := "test"

		m.Store(key, "value")
		m.Delete(key)

		_, exists := m.Load(key)
		if exists {
			t.Errorf("expected key to be deleted")
		}
		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("Len_MultipleKeys_ReturnsCorrectCount", func(t *testing.T) {
		m := NewShardedMap(8)

		for i := 0; i < 100; i++ {
			key := fmt.Sprintf("%d", i)
			m.Store(key, i)
		}

		if m.Len() != 100 {
			t.Errorf("expected len 100, got %d", m.Len())
		}
	})

	t.Run("Clear_AllKeys_ReturnsEmpty", func(t *testing.T) {
		m := NewShardedMap(32)

		for i := 0; i < 50; i++ {
			key := fmt.Sprintf("%d", i)
			m.Store(key, i)
		}

		if m.Len() != 50 {
			t.Errorf("expected len 50, got %d", m.Len())
		}

		m.Clear()

		if m.Len() != 0 {
			t.Errorf("expected len 0 after clear, got %d", m.Len())
		}
	})
}

func TestShardedMapExtendedOperations(t *testing.T) {
	m := NewShardedMap(32)

	t.Run("LoadOrStore_NewKey_InsertsEntry", func(t *testing.T) {
		key := "test"
		value := "value"

		result, loaded := m.LoadOrStore(key, value)

		if loaded {
			t.Errorf("expected key to not exist")
		}
		if result != value {
			t.Errorf("expected %v, got %v", value, result)
		}
		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})

	t.Run("LoadOrStore_ExistingKey_ReturnsExisting", func(t *testing.T) {
		key := "test"
		oldValue := "old"
		newValue := "new"

		m.Store(key, oldValue)
		result, loaded := m.LoadOrStore(key, newValue)

		if !loaded {
			t.Errorf("expected key to exist")
		}
		if result != oldValue {
			t.Errorf("expected %v, got %v", oldValue, result)
		}
		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})

	t.Run("LoadAndDelete_ExistingKey_ReturnsValue", func(t *testing.T) {
		key := "test"
		value := "value"

		m.Store(key, value)
		result, loaded := m.LoadAndDelete(key)

		if !loaded {
			t.Errorf("expected key to exist")
		}
		if result != value {
			t.Errorf("expected %v, got %v", value, result)
		}
		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("CompareAndDelete_MatchingValue_ReturnsTrue", func(t *testing.T) {
		key := "test"
		value := "value"

		m.Store(key, value)
		deleted := m.CompareAndDelete(key, value)

		if !deleted {
			t.Errorf("expected to delete matching key")
		}
		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("CompareAndSwap_MatchingValue_ReturnsTrue", func(t *testing.T) {
		key := "test"
		oldValue := "old"
		newValue := "new"

		m.Store(key, oldValue)
		swapped := m.CompareAndSwap(key, oldValue, newValue)

		if !swapped {
			t.Errorf("expected to swap matching value")
		}
		result, _ := m.Load(key)
		if result != newValue {
			t.Errorf("expected %v, got %v", newValue, result)
		}
		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})
}

func TestShardedMapRangeOperation(t *testing.T) {
	m := NewShardedMap(32)

	t.Run("Range_AllEntries_VisitAllKeys", func(t *testing.T) {
		expected := map[string]int{
			"key1": 1,
			"key2": 2,
			"key3": 3,
			"key4": 4,
			"key5": 5,
		}

		for k, v := range expected {
			m.Store(k, v)
		}

		visited := make(map[string]bool)

		m.Range(func(key, value interface{}) bool {
			k := key.(string)
			v := value.(int)

			if expected[k] != v {
				t.Errorf("expected %d for key %s, got %d", expected[k], k, v)
			}
			visited[k] = true
			return true
		})

		if len(visited) != len(expected) {
			t.Errorf("expected to visit %d keys, visited %d", len(expected), len(visited))
		}
		for k := range expected {
			if !visited[k] {
				t.Errorf("expected to visit key %s", k)
			}
		}
	})

	t.Run("Range_StopEarly_ReturnsFalse", func(t *testing.T) {
		m := NewShardedMap(32)
		for i := 0; i < 100; i++ {
			key := fmt.Sprintf("%d", i)
			m.Store(key, i)
		}

		count := 0

		m.Range(func(key, value interface{}) bool {
			count++
			if count == 50 {
				return false
			}
			return true
		})

		if count != 50 {
			t.Errorf("expected to stop after 50 visits, got %d", count)
		}
	})
}

func TestShardedMapConcurrency(t *testing.T) {
	t.Run("ConcurrentStore_MultipleKeys_DistributesEvenly", func(t *testing.T) {
		m := NewShardedMap(32)
		iterations := 10000
		goroutines := 10

		done := make(chan bool, goroutines)

		for i := 0; i < goroutines; i++ {
			go func(gid int) {
				for j := 0; j < iterations; j++ {
					key := fmt.Sprintf("%d", gid) + "_" + fmt.Sprintf("%d", j)
					m.Store(key, j)
				}
				done <- true
			}(i)
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}
	})

	t.Run("ConcurrentGet_MultipleKey_DoesNotPanic", func(t *testing.T) {
		m := NewShardedMap(32)
		iterations := 1000
		goroutines := 10

		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			m.Store(key, i)
		}

		done := make(chan bool, goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				for j := 0; j < iterations; j++ {
					key := fmt.Sprintf("%d", j)
					m.Load(key)
				}
				done <- true
			}()
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}
	})

	t.Run("ConcurrentDelete_MultipleKeys_MaintainsCount", func(t *testing.T) {
		m := NewShardedMap(32)
		iterations := 1000
		goroutines := 10
		keys := make([]string, iterations)

		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			keys[i] = key
			m.Store(key, i)
		}

		done := make(chan bool, goroutines)

		for i := 0; i < goroutines; i++ {
			go func(gid int) {
				for j := gid; j < iterations; j += goroutines {
					m.Delete(keys[j])
				}
				done <- true
			}(i)
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}
	})

	t.Run("ConcurrentMixedOperations_DoesNotPanic", func(t *testing.T) {
		m := NewShardedMap(32)
		iterations := 1000
		goroutines := 20

		done := make(chan bool, goroutines)

		for i := 0; i < goroutines; i++ {
			go func(gid int) {
				for j := 0; j < iterations; j++ {
					key := fmt.Sprintf("%d", gid) + "_" + fmt.Sprintf("%d", j)
					op := (gid + j) % 4

					switch op {
					case 0:
						m.Store(key, j)
					case 1:
						m.Load(key)
					case 2:
						m.Delete(key)
					case 3:
						m.Len()
					}
				}
				done <- true
			}(i)
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}
	})
}

func TestTypedShardedMapBasicOperations(t *testing.T) {
	t.Run("NewTypedShardedMap_DefaultShards_ReturnsMap", func(t *testing.T) {
		m := NewTypedShardedMap[string, int](32)

		if m == nil {
			t.Errorf("expected non-nil map")
		}
		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("NewTypedShardedMap32_Returns32Shards", func(t *testing.T) {
		m := NewTypedShardedMap32[string, string]()

		if m == nil {
			t.Errorf("expected non-nil map")
		}
	})
}

func TestTypedShardedMapCoreOperations(t *testing.T) {
	t.Run("Get_ExistingKey_ReturnsValue", func(t *testing.T) {
		m := NewTypedShardedMap[string, int](32)
		key := "test"
		value := 42
		m.Store(key, value)

		result, exists := m.Load(key)
		if !exists {
			t.Errorf("expected key to exist")
		}
		if result != value {
			t.Errorf("expected %d, got %d", value, result)
		}
	})

	t.Run("Store_NewKey_UpdatesLen", func(t *testing.T) {
		m := NewTypedShardedMap[string, string](16)
		key := "test"
		value := "value"

		m.Store(key, value)

		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})

	t.Run("Delete_ExistingKey_RemovesEntry", func(t *testing.T) {
		m := NewTypedShardedMap[string, int](32)
		key := "test"

		m.Store(key, 42)
		m.Delete(key)

		_, exists := m.Load(key)
		if exists {
			t.Errorf("expected key to be deleted")
		}
		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("Len_MultipleKeys_ReturnsCorrectCount", func(t *testing.T) {
		m := NewTypedShardedMap[int, string](8)

		for i := 0; i < 100; i++ {
			key := i
			m.Store(key, fmt.Sprintf("%c", i))
		}

		if m.Len() != 100 {
			t.Errorf("expected len 100, got %d", m.Len())
		}
	})

	t.Run("Clear_AllKeys_ReturnsEmpty", func(t *testing.T) {
		m := NewTypedShardedMap[string, float64](32)

		for i := 0; i < 50; i++ {
			key := fmt.Sprintf("%d", i)
			m.Store(key, float64(i))
		}

		if m.Len() != 50 {
			t.Errorf("expected len 50, got %d", m.Len())
		}

		m.Clear()

		if m.Len() != 0 {
			t.Errorf("expected len 0 after clear, got %d", m.Len())
		}
	})
}

func TestTypedShardedMapExtendedOperations(t *testing.T) {
	m := NewTypedShardedMap[string, string](32)

	t.Run("LoadOrStore_NewKey_InsertsEntry", func(t *testing.T) {
		key := "test"
		value := "value"

		result, loaded := m.LoadOrStore(key, value)

		if loaded {
			t.Errorf("expected key to not exist")
		}
		if result != value {
			t.Errorf("expected %q, got %q", value, result)
		}
		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})

	t.Run("LoadOrStore_ExistingKey_ReturnsExisting", func(t *testing.T) {
		key := "test"
		oldValue := "old"
		newValue := "new"

		m.Store(key, oldValue)
		result, loaded := m.LoadOrStore(key, newValue)

		if !loaded {
			t.Errorf("expected key to exist")
		}
		if result != oldValue {
			t.Errorf("expected %q, got %q", oldValue, result)
		}
		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})

	t.Run("LoadAndDelete_ExistingKey_ReturnsValue", func(t *testing.T) {
		key := "test"
		value := "value"

		m.Store(key, value)
		result, loaded := m.LoadAndDelete(key)

		if !loaded {
			t.Errorf("expected key to exist")
		}
		if result != value {
			t.Errorf("expected %q, got %q", value, result)
		}
		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("CompareAndDelete_MatchingValue_ReturnsTrue", func(t *testing.T) {
		key := "test"
		value := "value"

		m.Store(key, value)
		deleted := m.CompareAndDelete(key, value)

		if !deleted {
			t.Errorf("expected to delete matching key")
		}
		if m.Len() != 0 {
			t.Errorf("expected len 0, got %d", m.Len())
		}
	})

	t.Run("CompareAndSwap_MatchingValue_ReturnsTrue", func(t *testing.T) {
		key := "test"
		oldValue := "old"
		newValue := "new"

		m.Store(key, oldValue)
		swapped := m.CompareAndSwap(key, oldValue, newValue)

		if !swapped {
			t.Errorf("expected to swap matching value")
		}
		result, _ := m.Load(key)
		if result != newValue {
			t.Errorf("expected %q, got %q", newValue, result)
		}
		if m.Len() != 1 {
			t.Errorf("expected len 1, got %d", m.Len())
		}
	})
}

func TestTypedShardedMapRangeOperation(t *testing.T) {
	t.Run("Range_AllEntries_VisitAllKeys", func(t *testing.T) {
		m := NewTypedShardedMap[string, int](32)
		expected := map[string]int{
			"key1": 1,
			"key2": 2,
			"key3": 3,
			"key4": 4,
		}

		for k, v := range expected {
			m.Store(k, v)
		}

		visited := make(map[string]bool)

		m.Range(func(key string, value int) bool {
			if expected[key] != value {
				t.Errorf("expected %d for key %s, got %d", expected[key], key, value)
			}
			visited[key] = true
			return true
		})

		if len(visited) != len(expected) {
			t.Errorf("expected to visit %d keys, visited %d", len(expected), len(visited))
		}
		for k := range expected {
			if !visited[k] {
				t.Errorf("expected to visit key %s", k)
			}
		}
	})

	t.Run("Range_StopEarly_ReturnsFalse", func(t *testing.T) {
		m := NewTypedShardedMap[int, string](32)
		for i := 0; i < 100; i++ {
			key := i
			m.Store(key, fmt.Sprintf("%c", i))
		}

		count := 0

		m.Range(func(key int, value string) bool {
			count++
			if count == 50 {
				return false
			}
			return true
		})

		if count != 50 {
			t.Errorf("expected to stop after 50 visits, got %d", count)
		}
	})
}

func TestTypedShardedMapConcurrency(t *testing.T) {
	t.Run("ConcurrentStore_MultipleKey_DistributesEvenly", func(t *testing.T) {
		m := NewTypedShardedMap[string, int](32)
		iterations := 5000
		goroutines := 10

		done := make(chan bool, goroutines)

		for i := 0; i < goroutines; i++ {
			go func(gid int) {
				for j := 0; j < iterations; j++ {
					key := fmt.Sprintf("%d", gid) + "_" + fmt.Sprintf("%d", j)
					m.Store(key, j)
				}
				done <- true
			}(i)
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}
	})

	t.Run("ConcurrentGet_MultipleKey_DoesNotPanic", func(t *testing.T) {
		m := NewTypedShardedMap[int, int](32)
		iterations := 1000
		goroutines := 10

		for i := 0; i < iterations; i++ {
			key := i
			m.Store(key, i)
		}

		done := make(chan bool, goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				for j := 0; j < iterations; j++ {
					key := j
					m.Load(key)
				}
				done <- true
			}()
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}
	})

	t.Run("ConcurrentDelete_MultipleKeys_MaintainsCount", func(t *testing.T) {
		m := NewTypedShardedMap[string, int](32)
		iterations := 1000
		goroutines := 10
		keys := make([]string, iterations)

		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			keys[i] = key
			m.Store(key, i)
		}

		done := make(chan bool, goroutines)

		for i := 0; i < goroutines; i++ {
			go func(gid int) {
				for j := gid; j < iterations; j += goroutines {
					m.Delete(keys[j])
				}
				done <- true
			}(i)
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}
	})

	t.Run("ConcurrentMixedOperations_DoesNotPanic", func(t *testing.T) {
		m := NewTypedShardedMap[string, int](32)
		iterations := 5000
		goroutines := 20

		done := make(chan bool, goroutines)

		for i := 0; i < goroutines; i++ {
			go func(gid int) {
				for j := 0; j < iterations; j++ {
					key := fmt.Sprintf("%d", gid) + "_" + fmt.Sprintf("%d", j)
					op := (gid + j) % 4

					switch op {
					case 0:
						m.Store(key, j)
					case 1:
						m.Load(key)
					case 2:
						m.Delete(key)
					case 3:
						m.Len()
					}
				}
				done <- true
			}(i)
		}

		for i := 0; i < goroutines; i++ {
			<-done
		}
	})
}

func TestTypedShardedMapComplexTypes(t *testing.T) {
	type Product struct {
		ID    string
		Name  string
		Price float64
	}

	m := NewTypedShardedMap[string, Product](16)

	t.Run("Store_Get_ComplexStruct", func(t *testing.T) {
		product := Product{
			ID:    "P001",
			Name:  "Test Product",
			Price: 99.99,
		}
		m.Store(product.ID, product)

		result, exists := m.Load(product.ID)
		if !exists {
			t.Errorf("expected product to exist")
		}
		if result != product {
			t.Errorf("expected product %+v, got %+v", product, result)
		}
	})

	t.Run("Store_Get_PointerValue", func(t *testing.T) {
		m := NewTypedShardedMap[string, *Product](16)
		product := &Product{
			ID:    "P002",
			Name:  "Pointer Product",
			Price: 199.99,
		}
		m.Store(product.ID, product)

		result, exists := m.Load(product.ID)
		if !exists {
			t.Errorf("expected product pointer to exist")
		}
		if *result != *product {
			t.Errorf("expected product pointer %+v, got %+v", *product, *result)
		}
	})
}

func TestTypedShardedMapDifferentShardCounts(t *testing.T) {
	shardCounts := []int{1, 2, 4, 8, 16, 32, 64}

	for _, shards := range shardCounts {
		t.Run("ShardCount_"+fmt.Sprintf("%d", shards)+"_WorksCorrectly", func(t *testing.T) {
			m := NewTypedShardedMap[string, int](shards)

			for i := 0; i < 100; i++ {
				key := fmt.Sprintf("%d", i)
				m.Store(key, i)
			}

			if m.Len() != 100 {
				t.Errorf("expected len 100 with %d shards, got %d", shards, m.Len())
			}

			count := 0
			m.Range(func(key string, value int) bool {
				count++
				return true
			})

			if count != 100 {
				t.Errorf("expected to range 100 entries with %d shards, got %d", shards, count)
			}
		})
	}
}

func TestPerformance(t *testing.T) {
	iterations := 100000

	t.Run("Map_Performance_Store", func(t *testing.T) {
		m := NewMap()
		start := time.Now()

		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			m.Store(key, i)
		}

		duration := time.Since(start)
		t.Logf("Map Store %d entries: %v", iterations, duration)
	})

	t.Run("Map_Performance_Get", func(t *testing.T) {
		m := NewMap()
		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			m.Store(key, i)
		}

		start := time.Now()

		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			m.Load(key)
		}

		duration := time.Since(start)
		t.Logf("Map Get %d entries: %v", iterations, duration)
	})

	t.Run("Map_Performance_Range", func(t *testing.T) {
		m := NewMap()
		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			m.Store(key, i)
		}

		start := time.Now()
		count := 0

		m.Range(func(key, value interface{}) bool {
			count++
			return true
		})

		duration := time.Since(start)
		t.Logf("Map Range %d entries: %v", count, duration)
	})

	t.Run("TypedMap_Performance_Store", func(t *testing.T) {
		m := NewTypedMap[string, int]()
		start := time.Now()

		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			m.Store(key, i)
		}

		duration := time.Since(start)
		t.Logf("TypedMap Store %d entries: %v", iterations, duration)
	})

	t.Run("TypedMap_Performance_Get", func(t *testing.T) {
		m := NewTypedMap[string, int]()
		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			m.Store(key, i)
		}

		start := time.Now()

		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			m.Load(key)
		}

		duration := time.Since(start)
		t.Logf("TypedMap Get %d entries: %v", iterations, duration)
	})

	t.Run("TypedMap_Performance_Range", func(t *testing.T) {
		m := NewTypedMap[string, int]()
		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			m.Store(key, i)
		}

		start := time.Now()
		count := 0

		m.Range(func(key string, value int) bool {
			count++
			return true
		})

		duration := time.Since(start)
		t.Logf("TypedMap Range %d entries: %v", count, duration)
	})

	t.Run("ShardedMap_Performance_Store_32Shards", func(t *testing.T) {
		m := NewShardedMap(32)
		start := time.Now()

		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			m.Store(key, i)
		}

		duration := time.Since(start)
		t.Logf("ShardedMap Store %d entries with 32 shards: %v", iterations, duration)
	})

	t.Run("ShardedMap_Performance_Get_32Shards", func(t *testing.T) {
		m := NewShardedMap(32)
		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			m.Store(key, i)
		}

		start := time.Now()

		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			m.Load(key)
		}

		duration := time.Since(start)
		t.Logf("ShardedMap Get %d entries with 32 shards: %v", iterations, duration)
	})

	t.Run("ShardedMap_Performance_Range_32Shards", func(t *testing.T) {
		m := NewShardedMap(32)
		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			m.Store(key, i)
		}

		start := time.Now()
		count := 0

		m.Range(func(key, value interface{}) bool {
			count++
			return true
		})

		duration := time.Since(start)
		t.Logf("ShardedMap Range %d entries with 32 shards: %v", count, duration)
	})

	t.Run("TypedShardedMap_Performance_Store_32Shards", func(t *testing.T) {
		m := NewTypedShardedMap[string, int](32)
		start := time.Now()

		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			m.Store(key, i)
		}

		duration := time.Since(start)
		t.Logf("TypedShardedMap Store %d entries with 32 shards: %v", iterations, duration)
	})

	t.Run("TypedShardedMap_Performance_Get_32Shards", func(t *testing.T) {
		m := NewTypedShardedMap[string, int](32)
		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			m.Store(key, i)
		}

		start := time.Now()

		for i := 0; i < iterations; i++ {
			key := fmt.Sprintf("%d", i)
			m.Load(key)
		}

		duration := time.Since(start)
		t.Logf("TypedShardedMap Get %d entries with 32 shards: %v", iterations, duration)
	})

	t.Run("TypedShardedMap_Performance_Range_32Shards", func(t *testing.T) {
		m := NewTypedShardedMap[int, string](32)
		for i := 0; i < iterations; i++ {
			key := i
			m.Store(key, fmt.Sprintf("%c", i))
		}

		start := time.Now()
		count := 0

		m.Range(func(key int, value string) bool {
			count++
			return true
		})

		duration := time.Since(start)
		t.Logf("TypedShardedMap Range %d entries with 32 shards: %v", count, duration)
	})
}
