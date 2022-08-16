package pool

import (
	"context"
	"sync"

	"github.com/doron-cohen/poolher/worker"
)

type Pool[In interface{}, Out interface{}] struct {
	InChan  chan In
	OutChan chan worker.Result[Out]

	size     int
	workFunc worker.WorkFunc[In, Out]
	workers  []*worker.Worker[In, Out]
	wg       *sync.WaitGroup
}

func NewPool[In interface{}, Out interface{}](
	size int,
	workFunc worker.WorkFunc[In, Out],
) *Pool[In, Out] {
	return &Pool[In, Out]{
		InChan:   make(chan In, size),
		OutChan:  make(chan worker.Result[Out], size),
		size:     size,
		workFunc: workFunc,
		wg:       &sync.WaitGroup{},
	}
}

func (p *Pool[In, Out]) Start(ctx context.Context) context.CancelFunc {
	p.initWorkers()
	ctx, stop := context.WithCancel(ctx)
	p.wg.Add(p.size)

	go p.runWorkers(ctx)
	return stop
}

func (p *Pool[In, Out]) Wait() {
	p.wg.Wait()
}

func (p *Pool[In, Out]) initWorkers() {
	p.workers = make([]*worker.Worker[In, Out], p.size)
	for i := 0; i < p.size; i++ {
		p.workers[i] = worker.NewWorker(p.workFunc)
	}
}

func (p *Pool[In, Out]) runWorkers(ctx context.Context) {
	for _, w := range p.workers {
		go w.Run(ctx, p.InChan, p.OutChan, p.wg)
	}
}
