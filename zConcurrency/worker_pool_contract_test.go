package zConcurrency

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPoolSubmitRejectsFullQueueWithoutBlocking(t *testing.T) {
	pool := NewWorkerPool(1, 1)
	pool.Start()
	defer pool.Stop()

	started := make(chan struct{})
	release := make(chan struct{})
	if err := pool.Submit(func() error {
		close(started)
		<-release
		return nil
	}); err != nil {
		t.Fatalf("submit running task: %v", err)
	}
	<-started
	if err := pool.Submit(func() error { return nil }); err != nil {
		t.Fatalf("fill queue: %v", err)
	}

	start := time.Now()
	err := pool.Submit(func() error { return nil })
	if !errors.Is(err, ErrWorkerPoolFull) {
		t.Fatalf("full queue error = %v, want ErrWorkerPoolFull", err)
	}
	if elapsed := time.Since(start); elapsed > 100*time.Millisecond {
		t.Fatalf("non-blocking Submit took %v", elapsed)
	}
	close(release)
}

func TestWorkerPoolSubmitWithContextCancelsWhileFull(t *testing.T) {
	pool := NewWorkerPool(1, 1)
	pool.Start()
	defer pool.Stop()

	started := make(chan struct{})
	release := make(chan struct{})
	if err := pool.Submit(func() error {
		close(started)
		<-release
		return nil
	}); err != nil {
		t.Fatalf("submit running task: %v", err)
	}
	<-started
	if err := pool.Submit(func() error { return nil }); err != nil {
		t.Fatalf("fill queue: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()
	err := pool.SubmitWithContext(ctx, func() error { return nil })
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("SubmitWithContext error = %v, want deadline", err)
	}
	close(release)
	stats := pool.Stats()
	if stats.WaitCount != 1 || stats.WaitDuration <= 0 {
		t.Fatalf("wait stats = count %d duration %v, want one measured wait", stats.WaitCount, stats.WaitDuration)
	}
}

func TestWorkerPoolSubmitWithContextDoesNotCountImmediateCapacityAsWait(t *testing.T) {
	pool := NewWorkerPool(1, 1)
	pool.Start()
	defer pool.Stop()

	done := make(chan struct{})
	if err := pool.SubmitWithContext(context.Background(), func() error {
		close(done)
		return nil
	}); err != nil {
		t.Fatalf("submit with immediate capacity: %v", err)
	}
	<-done
	if stats := pool.Stats(); stats.WaitCount != 0 || stats.WaitDuration != 0 {
		t.Fatalf("immediate admission recorded as wait: count %d duration %v", stats.WaitCount, stats.WaitDuration)
	}
}

func TestWorkerPoolStopDrainsAcceptedAndRejectsConcurrentSubmit(t *testing.T) {
	pool := NewWorkerPool(4, 64)
	pool.Start()

	var executed atomic.Int64
	for i := 0; i < 64; i++ {
		if err := pool.SubmitWithContext(context.Background(), func() error {
			executed.Add(1)
			return nil
		}); err != nil {
			t.Fatalf("submit accepted task %d: %v", i, err)
		}
	}

	const submitters = 32
	var wg sync.WaitGroup
	errs := make(chan error, submitters)
	for i := 0; i < submitters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			errs <- pool.SubmitWithContext(ctx, func() error {
				executed.Add(1)
				return nil
			})
		}()
	}

	pool.Stop()
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil && !errors.Is(err, ErrWorkerPoolClosed) {
			t.Fatalf("concurrent submit error = %v", err)
		}
	}
	if got := executed.Load(); got < 64 {
		t.Fatalf("Stop did not drain accepted tasks: executed=%d", got)
	}
	if err := pool.Submit(func() error { return nil }); !errors.Is(err, ErrWorkerPoolClosed) {
		t.Fatalf("submit after Stop error = %v, want closed", err)
	}
}

func TestWorkerPoolSubmitReportsClosedWhileStopIsDraining(t *testing.T) {
	pool := NewWorkerPool(1, 1)
	pool.Start()

	started := make(chan struct{})
	release := make(chan struct{})
	if err := pool.Submit(func() error {
		close(started)
		<-release
		return nil
	}); err != nil {
		t.Fatalf("submit running task: %v", err)
	}
	<-started
	if err := pool.Submit(func() error { return nil }); err != nil {
		t.Fatalf("fill queue: %v", err)
	}

	stopped := make(chan struct{})
	go func() {
		pool.Stop()
		close(stopped)
	}()
	<-pool.stopCh
	if err := pool.Submit(func() error { return nil }); !errors.Is(err, ErrWorkerPoolClosed) {
		t.Fatalf("submit while Stop drains error = %v, want closed", err)
	}
	close(release)
	select {
	case <-stopped:
	case <-time.After(time.Second):
		t.Fatal("Stop did not finish draining")
	}
}
