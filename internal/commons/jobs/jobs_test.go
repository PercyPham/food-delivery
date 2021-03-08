package jobs_test

import (
	"context"
	"food-delivery/internal/commons/jobs"
	"testing"
)

func TestCreateNewJob(t *testing.T) {
	_ = jobs.New(func(ctx context.Context) error {
		return nil
	})
}
