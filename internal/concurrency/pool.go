package concurrency

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

// Task represents a unit of work that can be executed concurrently.
type Task func(ctx context.Context) error

// Pool implements a bounded worker pool with context cancellation,
// backpressure, error tracking, and graceful shutdown.
type Pool struct {
	maxWorkers int
	queue      chan Task
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	activeJobs int64
	closed     atomic.Bool
}

// NewPool initializes a new bounded worker pool.
func NewPool(ctx context.Context, maxWorkers int, queueCapacity int) *Pool {
	if maxWorkers <= 0 {
		maxWorkers = 4
	}
	if queueCapacity <= 0 {
		queueCapacity = maxWorkers * 2
	}

	pCtx, cancel := context.WithCancel(ctx)
	p := &Pool{
		maxWorkers: maxWorkers,
		queue:      make(chan Task, queueCapacity),
		ctx:        pCtx,
		cancel:     cancel,
	}

	// Spawn bounded workers
	for i := 0; i < maxWorkers; i++ {
		p.wg.Add(1)
		go p.worker()
	}

	return p
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case <-p.ctx.Done():
			return
		case task, ok := <-p.queue:
			if !ok {
				return
			}
			atomic.AddInt64(&p.activeJobs, 1)
			_ = task(p.ctx)
			atomic.AddInt64(&p.activeJobs, -1)
		}
	}
}

// Submit enqueues a task for processing. Returns an error if pool is closed or queue is full.
func (p *Pool) Submit(task Task) error {
	if p.closed.Load() {
		return fmt.Errorf("concurrency pool is closed")
	}
	if err := p.ctx.Err(); err != nil {
		return err
	}

	select {
	case <-p.ctx.Done():
		return p.ctx.Err()
	case p.queue <- task:
		return nil
	}
}

// ActiveJobs returns count of currently executing jobs.
func (p *Pool) ActiveJobs() int64 {
	return atomic.LoadInt64(&p.activeJobs)
}

// Stop initiates a graceful shutdown, draining remaining queued tasks.
func (p *Pool) Stop() {
	if p.closed.CompareAndSwap(false, true) {
		close(p.queue)
		p.wg.Wait()
		p.cancel()
	}
}
