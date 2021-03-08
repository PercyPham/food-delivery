package jobs_test

import (
	"context"
	"food-delivery/internal/commons/jobs"
	"testing"
)

func TestCreateNewJob(t *testing.T) {
	job := jobs.New(func(ctx context.Context) error {
		return nil
	})

	if job == nil {
		t.Errorf("Error creating new job")
	}
}
