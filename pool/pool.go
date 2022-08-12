package pool

import (
	"context"

	"github.com/doron-cohen/poolher/worker"
)

type Pool[In interface{}, Out interface{}] struct {
	InChan  chan In
	OutChan chan worker.Result[Out]

	workFunc worker.WorkFunc[In, Out]
	worker   *worker.Worker[In, Out]
}

func NewPool[In interface{}, Out interface{}](
	workFunc worker.WorkFunc[In, Out],
) *Pool[In, Out] {
	return &Pool[In, Out]{
		InChan:   make(chan In, 1),
		OutChan:  make(chan worker.Result[Out], 1),
		workFunc: workFunc,
	}
}

func (p *Pool[In, Out]) Start(ctx context.Context) context.CancelFunc {
	p.worker = worker.NewWorker(p.workFunc)
	ctx, stop := context.WithCancel(ctx)

	go p.worker.Run(ctx, p.InChan, p.OutChan)
	return stop
}
