package jobs

import (
	"context"
	"errors"
	"time"
)

// Job interface
type Job interface {
	Execute(ctx context.Context) error
	State() JobState
	HasRunCount() int
}

// JobState depicts state of a job
type JobState int

const (
	// StateInit states job has not been run
	StateInit JobState = iota
	// StateRunning states job is running
	StateRunning
	// StateFailed states job has failed
	StateFailed
	// StateFailedWithMaxRunCount states job has failed with max run count
	StateFailedWithMaxRunCount
	// StateCommpleted states job is completed
	StateCommpleted
)

var jobStateText = map[JobState]string{
	StateInit:                  "Init",
	StateRunning:               "Running",
	StateFailed:                "Failed",
	StateFailedWithMaxRunCount: "Failed With Max Run Count",
	StateCommpleted:            "Completed",
}

// JobStateText returns a text for the job state.
//   It returns the empty string if the code is unknown.
func JobStateText(state JobState) string { return jobStateText[state] }

type job struct {
	handler JobHandler

	state       JobState
	hasRunCount int
	maxRunCount int
	lastRunAt   time.Time
}

// JobHandler function structure
type JobHandler func(ctx context.Context) error

const (
	defaultMaxRunCount = 3
)

// New instance of Job
func New(h JobHandler) Job {
	return &job{
		handler:     h,
		state:       StateInit,
		hasRunCount: 0,
		maxRunCount: defaultMaxRunCount,
	}
}

func (j *job) Execute(ctx context.Context) error {
	if j.hasRunCount >= j.maxRunCount {
		j.state = StateFailedWithMaxRunCount
		return errors.New("Job has reached max run count") // TODO: change to apperror
	}

	j.state = StateRunning

	j.hasRunCount++

	err := j.handler(ctx)

	if err == nil {
		j.state = StateCommpleted
		return nil
	}

	if j.hasRunCount < j.maxRunCount {
		j.state = StateFailed
		return err
	}

	j.state = StateFailedWithMaxRunCount
	return err
}

func (j *job) State() JobState  { return j.state }
func (j *job) HasRunCount() int { return j.hasRunCount }
