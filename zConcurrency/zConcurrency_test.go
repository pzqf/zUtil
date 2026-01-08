package zConcurrency

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	// 创建一个工作池
	wp := NewWorkerPool(2, 10)
	wp.Start()
	defer wp.Stop()

	var wg sync.WaitGroup
	result := 0
	mu := sync.Mutex{}

	// 提交10个任务
	tasks := 10
	for i := 0; i < tasks; i++ {
		wg.Add(1)
		err := wp.Submit(func() error {
			defer wg.Done()
			
			// 模拟耗时操作
			time.Sleep(time.Millisecond * 10)
			
			mu.Lock()
			result++
			mu.Unlock()
			
			return nil
		})
		if err != nil {
			t.Errorf("Submit failed: %v", err)
		}
	}

	// 等待所有任务完成
	wg.Wait()

	// 验证结果
	if result != tasks {
		t.Errorf("expected %d, got %d", tasks, result)
	}
}

func TestWorkerPoolWithContext(t *testing.T) {
	// 创建一个工作池
	wp := NewWorkerPool(2, 10)
	wp.Start()
	defer wp.Stop()

	var wg sync.WaitGroup
	result := 0
	mu := sync.Mutex{}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 提交10个任务
	tasks := 10
	for i := 0; i < tasks; i++ {
		wg.Add(1)
		err := wp.SubmitWithContext(ctx, func() error {
			defer wg.Done()
			
			// 模拟耗时操作
			time.Sleep(time.Millisecond * 10)
			
			mu.Lock()
			result++
			mu.Unlock()
			
			return nil
		})
		if err != nil {
			t.Errorf("SubmitWithContext failed: %v", err)
		}
	}

	// 等待所有任务完成
	wg.Wait()

	// 验证结果
	if result != tasks {
		t.Errorf("expected %d, got %d", tasks, result)
	}
}

func TestSemaphore(t *testing.T) {
	// 创建一个信号量，允许最多2个并发
	sem := NewSemaphore(2)
	var wg sync.WaitGroup
	result := 0
	mu := sync.Mutex{}

	// 启动5个goroutine
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			// 获取信号量
			sem.Acquire(1)
			defer sem.Release(1)
			
			// 模拟耗时操作
			time.Sleep(time.Millisecond * 10)
			
			mu.Lock()
			result++
			mu.Unlock()
		}()
	}

	// 等待所有goroutine完成
	wg.Wait()

	// 验证结果
	if result != 5 {
		t.Errorf("expected %d, got %d", 5, result)
	}
}

func TestSemaphoreWithContext(t *testing.T) {
	// 创建一个信号量，允许最多1个并发
	sem := NewSemaphore(1)
	var wg sync.WaitGroup

	// 启动第一个goroutine，占用信号量
	wg.Add(1)
	go func() {
		defer wg.Done()
		sem.Acquire(1)
		// 占用较长时间
		time.Sleep(time.Millisecond * 100)
		sem.Release(1)
	}()

	// 等待第一个goroutine获取信号量
	time.Sleep(time.Millisecond * 10)

	// 尝试获取信号量，带超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 50)
	defer cancel()

	err := sem.AcquireWithContext(ctx, 1)
	if err == nil {
		t.Error("AcquireWithContext should have timed out")
	}

	// 等待所有goroutine完成
	wg.Wait()
}

func TestMutexWithTimeout(t *testing.T) {
	// 创建带超时的互斥锁
	mu := NewMutexWithTimeout()
	var wg sync.WaitGroup

	// 启动第一个goroutine，占用互斥锁
	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock()
		// 占用较长时间
		time.Sleep(time.Millisecond * 100)
		mu.Unlock()
	}()

	// 等待第一个goroutine获取互斥锁
	time.Sleep(time.Millisecond * 10)

	// 尝试获取互斥锁，带超时
	err := mu.LockWithTimeout(time.Millisecond * 50)
	if err == nil {
		t.Error("LockWithTimeout should have timed out")
	}

	// 等待所有goroutine完成
	wg.Wait()

	// 再次尝试获取互斥锁，这次应该成功
	err = mu.LockWithTimeout(time.Millisecond * 50)
	if err != nil {
		t.Errorf("LockWithTimeout failed: %v", err)
	}
	mu.Unlock()
}

func TestMutexWithContext(t *testing.T) {
	// 创建带超时的互斥锁
	mu := NewMutexWithTimeout()
	var wg sync.WaitGroup

	// 启动第一个goroutine，占用互斥锁
	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock()
		// 占用较长时间
		time.Sleep(time.Millisecond * 100)
		mu.Unlock()
	}()

	// 等待第一个goroutine获取互斥锁
	time.Sleep(time.Millisecond * 10)

	// 尝试获取互斥锁，带上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 50)
	defer cancel()

	err := mu.LockWithContext(ctx)
	if err == nil {
		t.Error("LockWithContext should have timed out")
	}

	// 等待所有goroutine完成
	wg.Wait()
}

func TestAtomicCounter(t *testing.T) {
	// 创建原子计数器
	counter := NewAtomicCounter(0)
	var wg sync.WaitGroup

	// 启动100个goroutine，每个递增100次
	goroutines := 100
	iterations := 100

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				counter.Increment()
			}
		}()
	}

	// 等待所有goroutine完成
	wg.Wait()

	// 验证结果
	expected := int64(goroutines * iterations)
	actual := counter.Get()
	if actual != expected {
		t.Errorf("expected %d, got %d", expected, actual)
	}

	// 测试其他方法
	counter.Set(100)
	if counter.Get() != 100 {
		t.Errorf("expected %d, got %d", 100, counter.Get())
	}

	counter.Add(50)
	if counter.Get() != 150 {
		t.Errorf("expected %d, got %d", 150, counter.Get())
	}

	counter.Decrement()
	if counter.Get() != 149 {
		t.Errorf("expected %d, got %d", 149, counter.Get())
	}

	// 测试比较并交换
	if !counter.CompareAndSwap(149, 200) {
		t.Error("CompareAndSwap should have succeeded")
	}
	if counter.Get() != 200 {
		t.Errorf("expected %d, got %d", 200, counter.Get())
	}
}

func TestRateLimiter(t *testing.T) {
	// 创建限流器，每秒允许10个请求，最大突发5个
	rl := NewRateLimiter(10, 5)

	// 测试突发请求
	allowed := 0
	for i := 0; i < 10; i++ {
		if rl.Allow() {
			allowed++
		}
	}

	// 应该允许5个突发请求
	if allowed != 5 {
		t.Errorf("expected 5 allowed requests, got %d", allowed)
	}

	// 等待令牌生成
	time.Sleep(time.Millisecond * 200)

	// 再次尝试请求
	allowed = 0
	for i := 0; i < 3; i++ {
		if rl.Allow() {
			allowed++
		}
	}

	// 应该允许2个请求（每秒10个，200ms约2个）
	if allowed != 2 {
		t.Errorf("expected 2 allowed requests, got %d", allowed)
	}
}

func TestRateLimiterAllowN(t *testing.T) {
	// 创建限流器，每秒允许10个请求，最大突发5个
	rl := NewRateLimiter(10, 5)

	// 测试一次请求多个令牌
	if !rl.AllowN(3) {
		t.Error("AllowN(3) should have succeeded")
	}

	// 再次请求3个令牌，应该失败（剩余2个令牌）
	if rl.AllowN(3) {
		t.Error("AllowN(3) should have failed")
	}
}

func TestWaitGroupWithContext(t *testing.T) {
	// 创建带上下文的WaitGroup
	wg := WaitGroupWithContext{}
	var mu sync.Mutex
	result := 0

	// 启动5个goroutine
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			// 模拟耗时操作
			time.Sleep(time.Millisecond * 10)
			
			mu.Lock()
			result++
			mu.Unlock()
		}()
	}

	// 等待所有goroutine完成
	errors := wg.Wait()

	// 验证结果
	if result != 5 {
		t.Errorf("expected %d, got %d", 5, result)
	}
	if len(errors) != 0 {
		t.Errorf("expected 0 errors, got %d", len(errors))
	}
}

func TestWaitGroupWithContextTimeout(t *testing.T) {
	// 创建带上下文的WaitGroup
	wg := WaitGroupWithContext{}

	// 启动一个长时间运行的goroutine
	wg.Add(1)
	go func() {
		// 占用较长时间
		time.Sleep(time.Millisecond * 100)
		wg.Done()
	}()

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 50)
	defer cancel()

	// 等待所有goroutine完成，应该超时
	_, err := wg.WaitWithContext(ctx)
	if err == nil {
		t.Error("WaitWithContext should have timed out")
	}
}