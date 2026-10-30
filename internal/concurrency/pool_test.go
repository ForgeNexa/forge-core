package concurrency

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestPool_BoundedExecutionAndGracefulDrain(t *testing.T) {
	ctx := context.Background()
	pool := NewPool(ctx, 4, 16)
	defer pool.Stop()

	var counter int64
	jobCount := 30

	for i := 0; i < jobCount; i++ {
		err := pool.Submit(func(c context.Context) error {
			time.Sleep(10 * time.Millisecond)
			atomic.AddInt64(&counter, 1)
			return nil
		})
		if err != nil {
			t.Fatalf("failed to submit job %d: %v", i, err)
		}
	}

	pool.Stop()

	if val := atomic.LoadInt64(&counter); val != int64(jobCount) {
		t.Fatalf("expected %d jobs completed, got %d", jobCount, val)
	}
}

func TestPool_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	pool := NewPool(ctx, 2, 4)

	cancel() // cancel immediately

	err := pool.Submit(func(c context.Context) error {
		return nil
	})

	if err == nil {
		t.Fatalf("expected error submitting to cancelled context pool, got nil")
	}

	pool.Stop()
}
