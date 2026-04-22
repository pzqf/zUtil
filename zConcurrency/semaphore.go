package zConcurrency

import (
	"context"
	"sync"
)

type Semaphore struct {
	permits int64
	mu      sync.Mutex
	ch      chan struct{}
}

func NewSemaphore(permits int64) *Semaphore {
	return &Semaphore{
		permits: permits,
		ch:      make(chan struct{}, 1),
	}
}

func (s *Semaphore) Acquire(n int64) {
	s.mu.Lock()
	for s.permits < n {
		s.mu.Unlock()
		<-s.ch
		s.mu.Lock()
	}
	s.permits -= n
	s.mu.Unlock()
}

func (s *Semaphore) AcquireWithContext(ctx context.Context, n int64) error {
	for {
		s.mu.Lock()
		if s.permits >= n {
			s.permits -= n
			s.mu.Unlock()
			return nil
		}
		s.mu.Unlock()

		select {
		case <-s.ch:
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *Semaphore) Release(n int64) {
	s.mu.Lock()
	s.permits += n
	s.mu.Unlock()

	select {
	case s.ch <- struct{}{}:
	default:
	}
}

func (s *Semaphore) AvailablePermits() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.permits
}
