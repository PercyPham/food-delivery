package jobs

import (
	"context"
	"sync"
)

// WorkerPool interface
type WorkerPool interface {
	Run(ctx context.Context) []error
}

type workerPool struct {
	Jobs []Job

	maxWorkers int
	collector  chan Job
	wg         sync.WaitGroup
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(jobs []Job, maxWorkers int) WorkerPool {
	return &workerPool{
		Jobs:       jobs,
		maxWorkers: maxWorkers,
		collector:  make(chan Job, len(jobs)),
	}
}

func (p *workerPool) Run(ctx context.Context) []error {
	var workerCount int
	if len(p.Jobs) < p.maxWorkers {
		workerCount = len(p.Jobs)
	} else {
		workerCount = p.maxWorkers
	}

	errChan := make(chan error, len(p.Jobs))

	for i := 0; i < workerCount; i++ {
		worker := NewWorker(p.collector, errChan)
		worker.Start(ctx, &p.wg)
	}

	for i := range p.Jobs {
		p.collector <- p.Jobs[i]
	}
	close(p.collector)
	p.wg.Wait()

	errs := make([]error, 0, len(p.Jobs))
	for i := 0; i < len(p.Jobs); i++ {
		if err := <-errChan; err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
