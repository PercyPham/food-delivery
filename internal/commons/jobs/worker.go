package jobs

import (
	"context"
	"time"
)

type worker struct {
}

func (w *worker) ExecJob(ctx context.Context, job Job) error {
	var err error

	for job.State() != StateFailedWithMaxRunCount {
		if job.State() != StateInit {
			durationUntilNextTry := calculateDurationUntilNextTry(job.HasRunCount())
			time.Sleep(durationUntilNextTry)
		}

		err = job.Execute(ctx)

		if job.State() == StateCommpleted {
			return nil
		}
	}

	return err
}

func calculateDurationUntilNextTry(hasRunCount int) time.Duration {
	switch hasRunCount {
	case 1:
		return time.Second
	default:
		return time.Duration(hasRunCount) * 5 * time.Second
	}
}
