package zConcurrency

import (
	"context"
	"sync"
	"time"
)

// MutexWithTimeout 带超时的互斥锁
type MutexWithTimeout struct {
	mu     sync.Mutex
	cond   *sync.Cond
	locked bool
}

// NewMutexWithTimeout 创建新的带超时互斥锁
func NewMutexWithTimeout() *MutexWithTimeout {
	m := &MutexWithTimeout{}
	m.cond = sync.NewCond(&m.mu)
	return m
}

// Lock 加锁（阻塞）
func (m *MutexWithTimeout) Lock() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for m.locked {
		m.cond.Wait()
	}
	m.locked = true
}

// TryLock 尝试加锁（非阻塞）
func (m *MutexWithTimeout) TryLock() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.locked {
		return false
	}
	m.locked = true
	return true
}

// LockWithTimeout 加锁（带超时）
func (m *MutexWithTimeout) LockWithTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return m.LockWithContext(ctx)
}

// LockWithContext 加锁（带上下文）
func (m *MutexWithTimeout) LockWithContext(ctx context.Context) error {
	for {
		m.mu.Lock()
		if !m.locked {
			m.locked = true
			m.mu.Unlock()
			return nil
		}

		done := make(chan struct{})
		go func() {
			m.mu.Lock()
			m.cond.Wait()
			m.mu.Unlock()
			close(done)
		}()

		m.mu.Unlock()

		select {
		case <-done:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Unlock 解锁
func (m *MutexWithTimeout) Unlock() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.locked {
		panic("unlock of unlocked mutex")
	}

	m.locked = false
	m.cond.Signal()
}
