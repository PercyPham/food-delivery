package jobs

import (
	"context"
	"sync"
	"time"
)

// Worker interface
type Worker interface {
	Start(ctx context.Context, wg *sync.WaitGroup)
}

type worker struct {
	jobChan chan Job
	errChan chan error
}

// NewWorker creates a new worker
func NewWorker(jobChan chan Job, errChan chan error) Worker {
	return &worker{jobChan, errChan}
}

func (w *worker) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for job := range w.jobChan {
			w.errChan <- handle(ctx, job)
		}
	}()
}

func handle(ctx context.Context, job Job) error {
	var err error

	for job.State() != StateFailedWithMaxRunCount {
		if job.State() != StateInit {
			durationUntilNextTry := calculateDurationUntilNextTry(job.RunCount())
			time.Sleep(durationUntilNextTry)
		}

		err = job.Execute(ctx)

		if job.State() == StateCommpleted {
			return nil
		}
	}

	return err
}

func calculateDurationUntilNextTry(runCount int) time.Duration {
	switch runCount {
	case 1:
		return time.Second
	default:
		return time.Duration(runCount) * 5 * time.Second
	}
}
