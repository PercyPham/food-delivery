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
	RunCount() int
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
	f JobHandler

	state       JobState
	runCount    int
	maxRunCount int
	lastRunAt   time.Time
}

// JobHandler function structure
type JobHandler func(ctx context.Context) error

const (
	defaultMaxRunCount = 3
)

// NewJob instance of Job
func NewJob(f JobHandler) Job {
	return &job{
		f:           f,
		state:       StateInit,
		runCount:    0,
		maxRunCount: defaultMaxRunCount,
	}
}

func (j *job) Execute(ctx context.Context) error {
	if j.runCount >= j.maxRunCount {
		j.state = StateFailedWithMaxRunCount
		return errors.New("Job has reached max run count") // TODO: change to apperror
	}

	j.state = StateRunning

	j.runCount++

	err := j.f(ctx)

	if err != nil {
		if j.runCount >= j.maxRunCount {
			j.state = StateFailedWithMaxRunCount
			return err
		}
		j.state = StateFailed
		return err
	}

	j.state = StateCommpleted
	return nil
}

func (j *job) State() JobState { return j.state }
func (j *job) RunCount() int   { return j.runCount }
