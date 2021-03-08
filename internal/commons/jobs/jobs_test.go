package jobs_test

import (
	"context"
	"errors"
	"fmt"
	"food-delivery/internal/commons/jobs"
	"testing"
	"time"
)

func TestJob(t *testing.T) {
	jobList := make([]jobs.Job, 0, 10)

	for i := 0; i < 10; i++ {
		idx := i
		jobList = append(jobList, jobs.NewJob(func(ctx context.Context) error {
			time.Sleep(time.Second)
			if idx%3 == 0 {
				return errors.New("err")
			}
			return nil
		}))
	}

	start := time.Now()

	workerPool := jobs.NewWorkerPool(jobList, 5)

	errs := workerPool.Run(context.Background())

	elapsed := time.Since(start)

	fmt.Println("Elapsed: ", elapsed)
	fmt.Println("Errs: ", errs)

}
