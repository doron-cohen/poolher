package worker

import (
	"context"
	"time"
)

type WorkFunc[In interface{}, Out interface{}] func(in In) (Out, error)
type Result[Out interface{}] struct {
	Value Out
	Error error
}

type Worker[In interface{}, Out interface{}] struct {
	workFunc WorkFunc[In, Out]
}

func NewWorker[In interface{}, Out interface{}](workFunc WorkFunc[In, Out]) *Worker[In, Out] {
	return &Worker[In, Out]{workFunc: workFunc}
}

func (w *Worker[In, Out]) Run(ctx context.Context, inChan chan In, outChan chan Result[Out]) {
MainLoop:
	for {
		select {
		case <-ctx.Done():
			break MainLoop
		case job := <-inChan:
			result, err := w.workFunc(job)
			outChan <- Result[Out]{Value: result, Error: err}
		}
		time.Sleep(time.Second)
	}
}
