package zConcurrency

import (
	"context"
	"sync"
)

// Semaphore 信号量结构体
type Semaphore struct {
	permits int64
	mu      sync.Mutex
	cond    *sync.Cond
}

// NewSemaphore 创建新的信号量
func NewSemaphore(permits int64) *Semaphore {
	s := &Semaphore{
		permits: permits,
	}
	s.cond = sync.NewCond(&s.mu)
	return s
}

// Acquire 获取n个信号量许可
func (s *Semaphore) Acquire(n int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for s.permits < n {
		s.cond.Wait()
	}
	s.permits -= n
}

// AcquireWithContext 获取n个信号量许可（带超时）
func (s *Semaphore) AcquireWithContext(ctx context.Context, n int64) error {
	for {
		s.mu.Lock()
		if s.permits >= n {
			s.permits -= n
			s.mu.Unlock()
			return nil
		}

		done := make(chan struct{})
		go func() {
			s.mu.Lock()
			s.cond.Wait()
			s.mu.Unlock()
			close(done)
		}()

		s.mu.Unlock()

		select {
		case <-done:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Release 释放n个信号量许可
func (s *Semaphore) Release(n int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.permits += n
	s.cond.Broadcast()
}

// AvailablePermits 获取当前可用的许可数
func (s *Semaphore) AvailablePermits() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.permits
}
