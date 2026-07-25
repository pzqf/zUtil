package zMap

import (
	"sync"
	"testing"
)

// TestStore_ConcurrentSameKeyCount 验证 OPT-13：多 goroutine 并发 Store 同一新键时，
// Len() 计数应恰为 1（此前检查-再动作会让多个 goroutine 各 +1 导致 Len 虚高）。
func TestStore_ConcurrentSameKeyCount(t *testing.T) {
	m := NewMap()
	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				m.Store("same-key", j)
			}
		}()
	}
	wg.Wait()
	if got := m.Len(); got != 1 {
		t.Fatalf("Len() = %d, want 1 (concurrent Store of same new key must not over-count)", got)
	}
}

// TestStore_ConcurrentDistinctKeysCount 并发 Store 不同键，Len() 应等于去重键数。
func TestStore_ConcurrentDistinctKeysCount(t *testing.T) {
	m := NewMap()
	const nKeys = 500
	var wg sync.WaitGroup
	// 每个键被多个 goroutine 重复写，最终计数应等于键数。
	for r := 0; r < 8; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := 0; k < nKeys; k++ {
				m.Store(k, k)
			}
		}()
	}
	wg.Wait()
	if got := m.Len(); got != int64(nKeys) {
		t.Fatalf("Len() = %d, want %d", got, nKeys)
	}
}
