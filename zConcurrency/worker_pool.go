package zConcurrency

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Task 表示一个可执行的任务。
type Task func() error

var (
	ErrWorkerPoolNotRunning = errors.New("worker pool is not running")
	ErrWorkerPoolFull       = errors.New("worker pool queue is full")
	ErrWorkerPoolClosed     = errors.New("worker pool is closed")
)

type workerPoolState uint8

const (
	workerPoolNew workerPoolState = iota
	workerPoolRunning
	workerPoolStopping
	workerPoolStopped
)

// WorkerPoolStats 是 admission 与排队等待的进程内快照。
type WorkerPoolStats struct {
	QueueCapacity int
	QueueDepth    int
	Accepted      uint64
	Rejected      uint64
	WaitCount     uint64
	WaitDuration  time.Duration
}

// WorkerPool 是一个显式启动、停止时排空已接纳任务的有界工作池。
//
// Submit 是非阻塞 admission；队列满时返回 ErrWorkerPoolFull。需要等待容量的调用者必须使用
// SubmitWithContext，让取消与 deadline 成为调用契约的一部分。
type WorkerPool struct {
	workers   int
	taskQueue chan Task
	slots     chan struct{}
	errors    chan error

	admissionMu sync.RWMutex
	state       workerPoolState
	stopCh      chan struct{}
	stoppedCh   chan struct{}
	stopOnce    sync.Once
	errorsOnce  sync.Once
	wg          sync.WaitGroup

	accepted     atomic.Uint64
	rejected     atomic.Uint64
	waitCount    atomic.Uint64
	waitDuration atomic.Int64
}

// NewWorkerPool 创建工作池。非法 worker/queue 参数收敛为 1，避免构造一个永远无法接纳任务的池。
func NewWorkerPool(workers int, queueSize int) *WorkerPool {
	if workers <= 0 {
		workers = 1
	}
	if queueSize <= 0 {
		queueSize = 1
	}

	slots := make(chan struct{}, queueSize)
	for i := 0; i < queueSize; i++ {
		slots <- struct{}{}
	}
	return &WorkerPool{
		workers:   workers,
		taskQueue: make(chan Task, queueSize),
		slots:     slots,
		errors:    make(chan error, workers),
		state:     workerPoolNew,
		stopCh:    make(chan struct{}),
		stoppedCh: make(chan struct{}),
	}
}

// Start 启动工作池。工作池是单次生命周期，Stop 后不会伪重启。
func (p *WorkerPool) Start() {
	p.admissionMu.Lock()
	if p.state != workerPoolNew {
		p.admissionMu.Unlock()
		return
	}
	p.state = workerPoolRunning
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
	p.admissionMu.Unlock()
}

func (p *WorkerPool) worker() {
	defer p.wg.Done()
	for {
		select {
		case task := <-p.taskQueue:
			p.releaseQueueSlot()
			p.runTask(task)
		case <-p.stopCh:
			// Stop 已拒绝新任务；此后 taskQueue 只减不增，可以稳定排空。
			for {
				select {
				case task := <-p.taskQueue:
					p.releaseQueueSlot()
					p.runTask(task)
				default:
					return
				}
			}
		}
	}
}

func (p *WorkerPool) runTask(task Task) {
	if task == nil {
		return
	}
	if err := task(); err != nil {
		select {
		case p.errors <- err:
		default:
		}
	}
}

func (p *WorkerPool) releaseQueueSlot() {
	p.slots <- struct{}{}
}

func (p *WorkerPool) stateError() error {
	switch p.state {
	case workerPoolNew:
		return ErrWorkerPoolNotRunning
	case workerPoolRunning:
		return nil
	default:
		return ErrWorkerPoolClosed
	}
}

// Submit 非阻塞地提交任务。
func (p *WorkerPool) Submit(task Task) error {
	select {
	case <-p.stopCh:
		p.rejected.Add(1)
		return ErrWorkerPoolClosed
	default:
	}
	select {
	case <-p.slots:
	default:
		p.rejected.Add(1)
		return ErrWorkerPoolFull
	}
	return p.commitTask(task)
}

// SubmitWithContext 等待队列容量，并在 ctx 取消、池关闭或任务接纳时返回。
func (p *WorkerPool) SubmitWithContext(ctx context.Context, task Task) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		p.rejected.Add(1)
		return err
	}
	select {
	case <-p.stopCh:
		p.rejected.Add(1)
		return ErrWorkerPoolClosed
	default:
	}
	select {
	case <-p.slots:
		return p.commitTask(task)
	default:
	}

	start := time.Now()
	p.waitCount.Add(1)
	select {
	case <-ctx.Done():
		p.waitDuration.Add(int64(time.Since(start)))
		p.rejected.Add(1)
		return ctx.Err()
	case <-p.stopCh:
		p.waitDuration.Add(int64(time.Since(start)))
		p.rejected.Add(1)
		return ErrWorkerPoolClosed
	case <-p.slots:
	}
	p.waitDuration.Add(int64(time.Since(start)))
	return p.commitTask(task)
}

// commitTask 在线性化 admission 的读锁内完成最终状态检查和入队。Stop 取得写锁后，后续提交不会再进入队列。
func (p *WorkerPool) commitTask(task Task) error {
	p.admissionMu.RLock()
	if err := p.stateError(); err != nil {
		p.admissionMu.RUnlock()
		p.releaseQueueSlot()
		p.rejected.Add(1)
		return err
	}
	p.taskQueue <- task // 已取得 slots 额度，必定不会阻塞。
	p.accepted.Add(1)
	p.admissionMu.RUnlock()
	return nil
}

// Stop 原子拒绝新任务，排空所有已接纳任务，再关闭错误流。并发/重复调用是幂等的。
func (p *WorkerPool) Stop() {
	p.stopOnce.Do(func() {
		p.admissionMu.Lock()
		switch p.state {
		case workerPoolNew:
			p.state = workerPoolStopped
			close(p.stopCh)
		case workerPoolRunning:
			p.state = workerPoolStopping
			close(p.stopCh)
		}
		p.admissionMu.Unlock()

		p.wg.Wait()
		p.admissionMu.Lock()
		p.state = workerPoolStopped
		p.admissionMu.Unlock()
		p.errorsOnce.Do(func() { close(p.errors) })
		close(p.stoppedCh)
	})
	<-p.stoppedCh
}

// Wait 等待当前 worker 退出；运行中的池只有在 Stop 后才会退出。
func (p *WorkerPool) Wait() {
	p.wg.Wait()
}

// Errors 获取执行过程中的错误流。
func (p *WorkerPool) Errors() <-chan error {
	return p.errors
}

// Stats 返回当前队列和 admission 指标快照。
func (p *WorkerPool) Stats() WorkerPoolStats {
	return WorkerPoolStats{
		QueueCapacity: cap(p.taskQueue),
		QueueDepth:    len(p.taskQueue),
		Accepted:      p.accepted.Load(),
		Rejected:      p.rejected.Load(),
		WaitCount:     p.waitCount.Load(),
		WaitDuration:  time.Duration(p.waitDuration.Load()),
	}
}
