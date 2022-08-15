package pool

import (
	"context"
	"testing"
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
