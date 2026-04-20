package zConcurrency

import (
	"sync/atomic"
)

// AtomicCounter 原子计数器
type AtomicCounter struct {
	value int64
}

// NewAtomicCounter 创建新的原子计数器
func NewAtomicCounter(initial int64) *AtomicCounter {
	return &AtomicCounter{
		value: initial,
	}
}

// Increment 原子递增
func (c *AtomicCounter) Increment() int64 {
	return atomic.AddInt64(&c.value, 1)
}

// Decrement 原子递减
func (c *AtomicCounter) Decrement() int64 {
	return atomic.AddInt64(&c.value, -1)
}

// Add 原子加法
func (c *AtomicCounter) Add(n int64) int64 {
	return atomic.AddInt64(&c.value, n)
}

// Get 获取当前值
func (c *AtomicCounter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}

// Set 设置当前值
func (c *AtomicCounter) Set(n int64) {
	atomic.StoreInt64(&c.value, n)
}

// CompareAndSwap 比较并交换
func (c *AtomicCounter) CompareAndSwap(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&c.value, old, new)
}
