package pool

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"distributed-scheduler/pkg/logger"
)

var (
	ErrPoolClosed = errors.New("协程池已关闭")
	ErrPoolFull   = errors.New("协程池已满")
	ErrTimeout    = errors.New("提交任务超时")
)

// Task 任务函数
type Task func()

// WorkerPool Goroutine池
type WorkerPool struct {
	maxWorkers  int32         // 最大工作协程数
	taskQueue   chan Task     // 任务队列
	workerCount int32         // 当前工作协程数
	running     int32         // 运行中的任务数
	closed      int32         // 是否已关闭
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewWorkerPool 创建协程池
func NewWorkerPool(maxWorkers, queueSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{
		maxWorkers: int32(maxWorkers),
		taskQueue:  make(chan Task, queueSize),
		ctx:        ctx,
		cancel:     cancel,
	}

	// 启动初始工作协程
	initialWorkers := maxWorkers / 4
	if initialWorkers < 1 {
		initialWorkers = 1
	}
	for i := 0; i < initialWorkers; i++ {
		pool.addWorker()
	}

	logger.Infof("协程池启动成功, 最大工作协程数: %d, 队列大小: %d", maxWorkers, queueSize)
	return pool
}

// addWorker 添加工作协程
func (p *WorkerPool) addWorker() {
	atomic.AddInt32(&p.workerCount, 1)
	p.wg.Add(1)

	go func() {
		defer func() {
			atomic.AddInt32(&p.workerCount, -1)
			p.wg.Done()
			if r := recover(); r != nil {
				logger.Errorf("协程池任务panic: %v", r)
			}
		}()

		idleTimeout := time.NewTimer(30 * time.Second)
		defer idleTimeout.Stop()

		for {
			select {
			case task, ok := <-p.taskQueue:
				if !ok {
					return
				}
				idleTimeout.Reset(30 * time.Second)

				atomic.AddInt32(&p.running, 1)
				func() {
					defer func() {
						atomic.AddInt32(&p.running, -1)
						if r := recover(); r != nil {
							logger.Errorf("任务执行panic: %v", r)
						}
					}()
					task()
				}()

			case <-idleTimeout.C:
				// 超时退出，但保证最少有1个worker
				if atomic.LoadInt32(&p.workerCount) > 1 {
					return
				}
				idleTimeout.Reset(30 * time.Second)

			case <-p.ctx.Done():
				return
			}
		}
	}()
}

// Submit 提交任务
func (p *WorkerPool) Submit(task Task) error {
	if atomic.LoadInt32(&p.closed) == 1 {
		return ErrPoolClosed
	}

	// 尝试直接放入队列
	select {
	case p.taskQueue <- task:
		// 检查是否需要扩容
		p.maybeExpand()
		return nil
	default:
		// 队列满了，尝试扩容
		if p.maybeExpand() {
			// 再次尝试
			select {
			case p.taskQueue <- task:
				return nil
			default:
				return ErrPoolFull
			}
		}
		return ErrPoolFull
	}
}

// SubmitWithTimeout 带超时提交任务
func (p *WorkerPool) SubmitWithTimeout(task Task, timeout time.Duration) error {
	if atomic.LoadInt32(&p.closed) == 1 {
		return ErrPoolClosed
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case p.taskQueue <- task:
		p.maybeExpand()
		return nil
	case <-timer.C:
		return ErrTimeout
	}
}

// maybeExpand 尝试扩容
func (p *WorkerPool) maybeExpand() bool {
	currentWorkers := atomic.LoadInt32(&p.workerCount)
	queueLen := len(p.taskQueue)

	// 如果队列积压较多且未达到最大工作协程数，则扩容
	if queueLen > int(currentWorkers) && currentWorkers < p.maxWorkers {
		p.addWorker()
		return true
	}
	return false
}

// Running 获取运行中的任务数
func (p *WorkerPool) Running() int {
	return int(atomic.LoadInt32(&p.running))
}

// WorkerCount 获取当前工作协程数
func (p *WorkerPool) WorkerCount() int {
	return int(atomic.LoadInt32(&p.workerCount))
}

// QueueSize 获取队列中等待的任务数
func (p *WorkerPool) QueueSize() int {
	return len(p.taskQueue)
}

// IsClosed 是否已关闭
func (p *WorkerPool) IsClosed() bool {
	return atomic.LoadInt32(&p.closed) == 1
}

// Close 关闭协程池
func (p *WorkerPool) Close() {
	if !atomic.CompareAndSwapInt32(&p.closed, 0, 1) {
		return
	}

	p.cancel()
	close(p.taskQueue)
	logger.Info("协程池开始关闭")
}

// Wait 等待所有任务完成
func (p *WorkerPool) Wait() {
	p.wg.Wait()
	logger.Info("协程池已完全关闭")
}

// Shutdown 优雅关闭
func (p *WorkerPool) Shutdown(timeout time.Duration) error {
	p.Close()

	done := make(chan struct{})
	go func() {
		p.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return errors.New("关闭超时")
	}
}

// Stats 获取统计信息
type Stats struct {
	MaxWorkers  int `json:"max_workers"`
	Workers     int `json:"workers"`
	Running     int `json:"running"`
	QueueSize   int `json:"queue_size"`
	QueueCap    int `json:"queue_cap"`
}

// GetStats 获取统计信息
func (p *WorkerPool) GetStats() Stats {
	return Stats{
		MaxWorkers: int(p.maxWorkers),
		Workers:    p.WorkerCount(),
		Running:    p.Running(),
		QueueSize:  p.QueueSize(),
		QueueCap:   cap(p.taskQueue),
	}
}

