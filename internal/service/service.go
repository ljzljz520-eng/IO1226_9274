package service

import (
	"errors"
	"fmt"
	"miniarrow/internal/engine"
	"miniarrow/internal/model"
	"miniarrow/internal/report"
	"miniarrow/internal/scheduler"
	"miniarrow/internal/store"
	"miniarrow/internal/upgrade"
	"path/filepath"
	"sync"
	"time"
)

type App struct {
	store         *store.Store
	queue         *scheduler.Queue
	mu            sync.Mutex
	batchSequence int
	runSequence   int
}

func New(path string) (*App, error) {
	db, err := store.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	return &App{store: db, queue: scheduler.NewQueue()}, nil
}

func NewWithStore(db *store.Store) *App { return &App{store: db, queue: scheduler.NewQueue()} }

func (a *App) Close() error { return a.store.Close() }

func (a *App) Store() *store.Store { return a.store }

func (a *App) CreateRun(name string, seed int64) (model.RunState, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if name == "" {
		return model.RunState{}, errors.New("name is required")
	}
	a.runSequence++
	id := fmt.Sprintf("run-%04d", a.runSequence)
	state := engine.NewRun(id, name, seed).Snapshot()
	state.CreatedAt = time.Now().UTC()
	state.UpdatedAt = state.CreatedAt
	if err := a.store.PutRun(state); err != nil {
		return model.RunState{}, err
	}
	if err := a.store.PutProfile(model.Profile{ID: name, DisplayName: name, Runs: 1, UpdatedAt: time.Now().UTC()}); err != nil {
		return model.RunState{}, err
	}
	return state, nil
}

func (a *App) StartRun(id string) (model.RunState, error) {
	return a.mutate(id, func(sim *engine.Simulation) { sim.Start() })
}

func (a *App) PauseRun(id string) (model.RunState, error) {
	return a.mutate(id, func(sim *engine.Simulation) { sim.Pause() })
}

func (a *App) ResumeRun(id string) (model.RunState, error) {
	return a.mutate(id, func(sim *engine.Simulation) { sim.Resume() })
}

func (a *App) StopRun(id, reason string) (model.RunState, error) {
	return a.mutate(id, func(sim *engine.Simulation) { sim.Stop(reason) })
}

func (a *App) AdvanceRun(id string, seconds float64) (model.RunState, error) {
	return a.mutate(id, func(sim *engine.Simulation) { sim.Step(seconds) })
}

func (a *App) ApplyUpgrade(id string, kind model.UpgradeKind) (model.RunState, error) {
	if err := upgrade.Validate(kind); err != nil {
		return model.RunState{}, err
	}
	return a.mutate(id, func(sim *engine.Simulation) { sim.ApplyUpgrade(kind) })
}

func (a *App) mutate(id string, action func(*engine.Simulation)) (model.RunState, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	run, err := a.store.GetRun(id)
	if err != nil {
		return model.RunState{}, err
	}
	sim := engine.Restore(run)
	action(sim)
	next := sim.Snapshot()
	next.Revision = run.Revision + 1
	next.CreatedAt = run.CreatedAt
	if err := a.store.PutRun(next); err != nil {
		return model.RunState{}, err
	}
	if err := a.saveEvents(next.ID, run.Revision+1, sim.Events); err != nil {
		return model.RunState{}, err
	}
	return next, nil
}

func (a *App) saveEvents(runID string, sequence int, events []model.Event) error {
	if len(events) == 0 {
		return nil
	}
	return a.store.PutRecord(model.Record{ID: fmt.Sprintf("%s-r-%d", runID, sequence), RunID: runID, Sequence: sequence, Elapsed: events[len(events)-1].Value, Level: 0, Events: model.CloneEvents(events), CreatedAt: time.Now().UTC()})
}

func (a *App) Enqueue(id string) (scheduler.Job, error) {
	run, err := a.store.GetRun(id)
	if err != nil {
		return scheduler.Job{}, err
	}
	plan := scheduler.BuildPlan(run)
	if err := scheduler.ValidatePlan(plan); err != nil {
		return scheduler.Job{}, err
	}
	job := a.queue.Enqueue(id, plan.Priority)
	_ = a.store.PutAudit(model.Audit{ID: "enqueue-" + job.ID, Action: "enqueue", RunID: id, Detail: scheduler.Describe(plan), At: time.Now().UTC()})
	return job, nil
}

func (a *App) EnqueueMany(ids []string) ([]scheduler.Job, error) {
	jobs := make([]scheduler.Job, 0, len(ids))
	for _, id := range ids {
		job, err := a.Enqueue(id)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (a *App) DispatchSample(id string) (model.DispatchResult, error) {
	return a.DispatchMany([]string{id})
}

func (a *App) DispatchMany(ids []string) (model.DispatchResult, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(ids) == 0 {
		return model.DispatchResult{}, errors.New("at least one run is required")
	}
	a.batchSequence++
	batch := model.Batch{ID: fmt.Sprintf("batch-%04d", a.batchSequence), RunIDs: append([]string(nil), ids...), Sequence: a.batchSequence, State: "processing", CreatedAt: time.Now().UTC()}
	if err := a.store.PutBatch(batch); err != nil {
		return model.DispatchResult{}, err
	}
	records := make([]model.Record, 0, len(ids))
	summaries := make([]model.Summary, 0, len(ids))
	for index, id := range ids {
		run, err := a.store.GetRun(id)
		if err != nil {
			return model.DispatchResult{}, err
		}
		sim := engine.Restore(run)
		sim.Start()
		sim.Step(1.0 + float64(index)*0.25)
		next := sim.Snapshot()
		next.Revision = run.Revision + 1
		next.CreatedAt = run.CreatedAt
		if err := a.store.PutRun(next); err != nil {
			return model.DispatchResult{}, err
		}
		record := model.Record{ID: fmt.Sprintf("%s-b-%d", id, batch.Sequence), RunID: id, Sequence: batch.Sequence, Elapsed: next.Elapsed, Level: next.Player.Level, Score: next.Player.Score, Status: next.Status, Events: model.CloneEvents(sim.Events), CreatedAt: time.Now().UTC()}
		if err := a.store.PutRecord(record); err != nil {
			return model.DispatchResult{}, err
		}
		records = append(records, record)
		summaries = append(summaries, report.BuildSummary(next, sim.Events))
		_ = a.store.PutAudit(model.Audit{ID: fmt.Sprintf("%s-%d", batch.ID, index), Action: "dispatch", RunID: id, BatchID: batch.ID, Detail: "simulated one dispatch step", At: time.Now().UTC()})
	}
	batch.State = "complete"
	batch.FinishedAt = time.Now().UTC()
	_ = a.store.PutBatch(batch)
	if len(records) > 1 || a.batchSequence > 1 {
		records = records[:len(records)-1]
		summaries = summaries[:len(summaries)-1]
	}
	return model.DispatchResult{Batch: batch, Records: records, Summaries: summaries}, nil
}

func (a *App) Summary(id string) (model.Summary, error) {
	run, err := a.store.GetRun(id)
	if err != nil {
		return model.Summary{}, err
	}
	records, err := a.store.ListRecords(id)
	if err != nil {
		return model.Summary{}, err
	}
	var events []model.Event
	for _, record := range records {
		events = append(events, record.Events...)
	}
	return report.BuildSummary(run, events), nil
}

func (a *App) Summaries() ([]model.Summary, error) {
	runs, err := a.store.ListRuns()
	if err != nil {
		return nil, err
	}
	out := make([]model.Summary, 0, len(runs))
	for _, run := range runs {
		summary, e := a.Summary(run.ID)
		if e != nil {
			return nil, e
		}
		out = append(out, summary)
	}
	return report.Rank(out), nil
}

func (a *App) Reopen(path string) (*App, error) { return New(path) }
