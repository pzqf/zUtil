package zConcurrency

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Task 表示一个可执行的任务
type Task func() error

// WorkerPool 工作池结构体
type WorkerPool struct {
	workers   int
	taskQueue chan Task
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	errors    chan error
	running   atomic.Bool
}

// NewWorkerPool 创建新的工作池
func NewWorkerPool(workers int, queueSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers:   workers,
		taskQueue: make(chan Task, queueSize),
		wg:        sync.WaitGroup{},
		ctx:       ctx,
		cancel:    cancel,
		errors:    make(chan error, workers),
		running:   atomic.Bool{},
	}
}

// Start 启动工作池
func (p *WorkerPool) Start() {
	if !p.running.CompareAndSwap(false, true) {
		return
	}

	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

// worker 工作协程
func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case task, ok := <-p.taskQueue:
			if !ok {
				return
			}
			if err := task(); err != nil {
				select {
				case p.errors <- err:
				default:
					// 错误通道已满，忽略
				}
			}
		}
	}
}

// Submit 提交任务到工作池
func (p *WorkerPool) Submit(task Task) error {
	select {
	case <-p.ctx.Done():
		return fmt.Errorf("worker pool is closed")
	case p.taskQueue <- task:
		return nil
	}
}

// SubmitWithContext 提交任务到工作池（带上下文）
func (p *WorkerPool) SubmitWithContext(ctx context.Context, task Task) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-p.ctx.Done():
		return fmt.Errorf("worker pool is closed")
	case p.taskQueue <- task:
		return nil
	}
}

// Stop 停止工作池
func (p *WorkerPool) Stop() {
	if !p.running.CompareAndSwap(true, false) {
		return
	}
	p.cancel()
	close(p.taskQueue)
	p.wg.Wait()
	close(p.errors)
}

// Wait 等待所有任务完成
func (p *WorkerPool) Wait() {
	p.wg.Wait()
}

// Errors 获取执行过程中的错误
func (p *WorkerPool) Errors() <-chan error {
	return p.errors
}

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
		// 尝试获取许可
		s.mu.Lock()
		if s.permits >= n {
			s.permits -= n
			s.mu.Unlock()
			return nil
		}

		// 创建用于等待的channel
		done := make(chan struct{})
		go func() {
			// 在goroutine中再次获取锁并等待
			s.mu.Lock()
			s.cond.Wait()
			s.mu.Unlock()
			close(done)
		}()

		s.mu.Unlock()

		// 等待条件满足或上下文取消
		select {
		case <-done:
			// 条件满足，再次尝试
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

// MutexWithTimeout 带超时的互斥锁
type MutexWithTimeout struct {
	mu    sync.Mutex
	cond  *sync.Cond
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
		// 尝试加锁
		m.mu.Lock()
		if !m.locked {
			m.locked = true
			m.mu.Unlock()
			return nil
		}

		// 创建用于等待的channel
		done := make(chan struct{})
		go func() {
			// 在goroutine中再次获取锁并等待
			m.mu.Lock()
			m.cond.Wait()
			m.mu.Unlock()
			close(done)
		}()

		m.mu.Unlock()

		// 等待条件满足或上下文取消
		select {
		case <-done:
			// 条件满足，再次尝试
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

// RateLimiter 限流器结构体
type RateLimiter struct {
	rate     int64         // 每秒允许的请求数
	burst    int64         // 最大突发请求数
	tokens   int64         // 当前可用令牌数
	lastTime atomic.Int64  // 上次更新时间戳（纳秒）
	mu       sync.Mutex    // 保护tokens的互斥锁
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(rate, burst int64) *RateLimiter {
	if rate <= 0 {
		rate = 1
	}
	if burst <= 0 {
		burst = 1
	}
	rl := &RateLimiter{
		rate:     rate,
		burst:    burst,
		tokens:   burst,
		lastTime: atomic.Int64{},
	}
	// 初始化lastTime为当前时间，避免第一次调用时elapsed值过大导致整数溢出
	rl.lastTime.Store(time.Now().UnixNano())
	return rl
}

// Allow 是否允许请求通过
func (rl *RateLimiter) Allow() bool {
	return rl.AllowN(1)
}

// AllowN 是否允许n个请求通过
func (rl *RateLimiter) AllowN(n int64) bool {
	if n <= 0 {
		return true
	}

	now := time.Now().UnixNano()
	lastTime := rl.lastTime.Load()
	elapsed := now - lastTime

	// 更新令牌数
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// 再次检查时间，避免并发问题
	now = time.Now().UnixNano()
	elapsed = now - rl.lastTime.Load()

	// 计算新生成的令牌数
	tokensToAdd := elapsed * rl.rate / 1e9
	if tokensToAdd > 0 {
		// 更新最后时间
		rl.lastTime.Store(now)
		// 增加令牌数，但不超过burst
		rl.tokens = min(rl.tokens+tokensToAdd, rl.burst)
	}

	// 检查是否有足够的令牌
	if rl.tokens >= n {
		rl.tokens -= n
		return true
	}

	return false
}

// Wait 等待直到允许请求通过
func (rl *RateLimiter) Wait() {
	for !rl.Allow() {
		time.Sleep(time.Millisecond * 1)
	}
}

// WaitWithContext 等待直到允许请求通过（带上下文）
func (rl *RateLimiter) WaitWithContext(ctx context.Context) error {
	for {
		if rl.Allow() {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Millisecond * 1):
			// 继续尝试
		}
	}
}

// Reset 重置限流器
func (rl *RateLimiter) Reset() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.tokens = rl.burst
	rl.lastTime.Store(time.Now().UnixNano())
}

// min 返回两个int64的最小值
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// WaitGroupWithContext 带上下文的WaitGroup
type WaitGroupWithContext struct {
	wg sync.WaitGroup
	mu sync.Mutex
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