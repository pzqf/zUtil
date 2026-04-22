package zConcurrency

import (
	"context"
	"sync"
	"time"
)

type MutexWithTimeout struct {
	mu     sync.Mutex
	locked bool
	ch     chan struct{}
}

func NewMutexWithTimeout() *MutexWithTimeout {
	return &MutexWithTimeout{
		ch: make(chan struct{}, 1),
	}
}

func (m *MutexWithTimeout) Lock() {
	m.mu.Lock()
	if !m.locked {
		m.locked = true
		m.mu.Unlock()
		return
	}
	m.mu.Unlock()
	<-m.ch
	m.mu.Lock()
	m.locked = true
	m.mu.Unlock()
}

func (m *MutexWithTimeout) TryLock() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.locked {
		return false
	}
	m.locked = true
	return true
}

func (m *MutexWithTimeout) LockWithTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return m.LockWithContext(ctx)
}

func (m *MutexWithTimeout) LockWithContext(ctx context.Context) error {
	m.mu.Lock()
	if !m.locked {
		m.locked = true
		m.mu.Unlock()
		return nil
	}
	m.mu.Unlock()

	select {
	case <-m.ch:
		m.mu.Lock()
		m.locked = true
		m.mu.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *MutexWithTimeout) Unlock() {
	m.mu.Lock()
	if !m.locked {
		m.mu.Unlock()
		panic("unlock of unlocked mutex")
	}
	m.locked = false
	m.mu.Unlock()

	select {
	case m.ch <- struct{}{}:
	default:
	}
}
