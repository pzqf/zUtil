package zConcurrency

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
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
