package zConcurrency

import (
	"context"
	"sync"
)

// WaitGroupWithContext 带上下文的WaitGroup
type WaitGroupWithContext struct {
	wg     sync.WaitGroup
	mu     sync.Mutex
	errors []error
}

// Add 添加任务数
func (w *WaitGroupWithContext) Add(delta int) {
	w.wg.Add(delta)
}

// Done 标记任务完成
func (w *WaitGroupWithContext) Done() {
	w.wg.Done()
}

// Wait 等待所有任务完成
func (w *WaitGroupWithContext) Wait() []error {
	w.wg.Wait()
	return w.errors
}

// WaitWithContext 等待所有任务完成（带上下文）
func (w *WaitGroupWithContext) WaitWithContext(ctx context.Context) ([]error, error) {
	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return w.errors, ctx.Err()
	case <-done:
		return w.errors, nil
	}
}

// RecordError 记录错误
func (w *WaitGroupWithContext) RecordError(err error) {
	if err == nil {
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()
	w.errors = append(w.errors, err)
}
