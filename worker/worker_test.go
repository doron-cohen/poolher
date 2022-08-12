package worker

import (
	"context"
	"testing"
)

func TestWorkerJob(t *testing.T) {
	worker := NewWorker(func(in int) (int, error) {
		return in * 2, nil
	})

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	inChan := make(chan int, 1)
	outChan := make(chan Result[int], 1)
	go worker.Run(ctx, inChan, outChan)

	inChan <- 2
	result := <-outChan
	cancel()
	if result.Error != nil {
		t.Fatalf("Unexpected error returned: %s", result.Error)
	}
	if result.Value != 4 {
		t.Fatalf("Wrong value returned: %d", result.Value)
	}
}
