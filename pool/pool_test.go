package pool

import (
	"context"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	pool := NewPool(1, func(in int) (int, error) {
		return in * 2, nil
	})

	ctx := context.Background()
	stop := pool.Start(ctx)

	pool.InChan <- 2
	result := <-pool.OutChan
	stop()
	if result.Error != nil {
		t.Fatalf("Unexpected error returned: %s", result.Error)
	}
	if result.Value != 4 {
		t.Fatalf("Wrong value returned: %d", result.Value)
	}
}

func TestPoolMultipleWorkers(t *testing.T) {
	pool := NewPool(4, func(in int) (int, error) {
		return in * 2, nil
	})

	ctx := context.Background()
	stop := pool.Start(ctx)

	for i := 0; i < 4; i++ {
		pool.InChan <- 2
	}

	for i := 0; i < 4; i++ {
		result := <-pool.OutChan
		if result.Error != nil {
			t.Fatalf("Unexpected error returned: %s", result.Error)
		}
		if result.Value != 4 {
			t.Fatalf("Wrong value returned: %d", result.Value)
		}
	}

	stop()
}

func TestPoolWait(t *testing.T) {
	pool := NewPool(4, func(in int) (int, error) {
		time.Sleep(time.Second)
		return 1, nil
	})

	ctx := context.Background()
	stop := pool.Start(ctx)

	for i := 0; i < 4; i++ {
		pool.InChan <- 1
	}

	stop()
	stoppedAt := time.Now()

	pool.Wait()
	timeWaited := time.Since(stoppedAt)
	if timeWaited < time.Second {
		t.Fatalf("Wait didn't wait for pool. Only %s passed", timeWaited)
	}
}

func BenchmarkPool(b *testing.B) {
	pool := NewPool(32, func(in int) (int, error) {
		return in * 2, nil
	})

	ctx := context.Background()
	stop := pool.Start(ctx)

	go func() {
		// Drain outchan until stop is called
	LOOP:
		for {
			select {
			case <-ctx.Done():
				break LOOP
			case <-pool.OutChan:
				continue
			}
		}
	}()

	for i := 0; i < 10_000; i++ {
		pool.InChan <- 2
	}

	stop()
}
