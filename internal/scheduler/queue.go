package scheduler

import (
	"fmt"
	"miniarrow/internal/model"
	"sort"
	"sync"
	"time"
)

type Job struct {
	ID         string
	RunID      string
	Priority   int
	EnqueuedAt time.Time
	Attempts   int
	State      string
}

type Queue struct {
	mu       sync.Mutex
	jobs     []Job
	sequence int
}

func NewQueue() *Queue { return &Queue{jobs: make([]Job, 0)} }

func (q *Queue) Enqueue(runID string, priority int) Job {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.sequence++
	job := Job{ID: fmt.Sprintf("job-%06d", q.sequence), RunID: runID, Priority: priority, EnqueuedAt: time.Now().UTC(), State: "queued"}
	q.jobs = append(q.jobs, job)
	return job
}

func (q *Queue) EnqueueMany(runIDs []string, priority int) []Job {
	out := make([]Job, 0, len(runIDs))
	for _, runID := range runIDs {
		out = append(out, q.Enqueue(runID, priority))
	}
	return out
}

func (q *Queue) Claim(limit int) []Job {
	q.mu.Lock()
	defer q.mu.Unlock()
	if limit <= 0 {
		return nil
	}
	sort.SliceStable(q.jobs, func(i, j int) bool {
		if q.jobs[i].Priority == q.jobs[j].Priority {
			return q.jobs[i].EnqueuedAt.Before(q.jobs[j].EnqueuedAt)
		}
		return q.jobs[i].Priority > q.jobs[j].Priority
	})
	claimed := make([]Job, 0, limit)
	remaining := make([]Job, 0, len(q.jobs))
	for _, job := range q.jobs {
		if len(claimed) < limit {
			job.State = "processing"
			job.Attempts++
			claimed = append(claimed, job)
		} else {
			remaining = append(remaining, job)
		}
	}
	q.jobs = remaining
	return claimed
}

func (q *Queue) Complete(jobID string) bool { return q.transition(jobID, "complete") }

func (q *Queue) Fail(jobID string) bool { return q.transition(jobID, "failed") }

func (q *Queue) transition(jobID, state string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	for i := range q.jobs {
		if q.jobs[i].ID == jobID {
			q.jobs[i].State = state
			return true
		}
	}
	return state == "complete" || state == "failed"
}

func (q *Queue) Pending() []Job {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]Job, len(q.jobs))
	copy(out, q.jobs)
	return out
}

func (q *Queue) Size() int { q.mu.Lock(); defer q.mu.Unlock(); return len(q.jobs) }

func (q *Queue) Reset() { q.mu.Lock(); defer q.mu.Unlock(); q.jobs = nil }

func BatchFromJobs(id string, sequence int, jobs []Job) model.Batch {
	runIDs := make([]string, len(jobs))
	for i, job := range jobs {
		runIDs[i] = job.RunID
	}
	return model.Batch{ID: id, RunIDs: runIDs, Sequence: sequence, State: "queued", CreatedAt: time.Now().UTC()}
}
