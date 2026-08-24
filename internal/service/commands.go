package service

import (
	"errors"
	"miniarrow/internal/engine"
	"miniarrow/internal/model"
	"miniarrow/internal/report"
	"miniarrow/internal/upgrade"
	"sort"
)

type Command struct {
	Name    string
	RunID   string
	Seconds float64
	Upgrade model.UpgradeKind
	Reason  string
}

type CommandResult struct {
	State   model.RunState
	Summary model.Summary
	Message string
}

func (a *App) Execute(command Command) (CommandResult, error) {
	if command.RunID == "" {
		return CommandResult{}, errors.New("run id is required")
	}
	switch command.Name {
	case "start":
		state, err := a.StartRun(command.RunID)
		return CommandResult{State: state, Message: "started"}, err
	case "pause":
		state, err := a.PauseRun(command.RunID)
		return CommandResult{State: state, Message: "paused"}, err
	case "resume":
		state, err := a.ResumeRun(command.RunID)
		return CommandResult{State: state, Message: "resumed"}, err
	case "step":
		state, err := a.AdvanceRun(command.RunID, command.Seconds)
		return CommandResult{State: state, Message: "advanced"}, err
	case "upgrade":
		state, err := a.ApplyUpgrade(command.RunID, command.Upgrade)
		return CommandResult{State: state, Message: "upgrade applied"}, err
	case "stop":
		state, err := a.StopRun(command.RunID, command.Reason)
		return CommandResult{State: state, Message: "stopped"}, err
	case "summary":
		summary, err := a.Summary(command.RunID)
		return CommandResult{Summary: summary, Message: report.Format(summary)}, err
	default:
		return CommandResult{}, errors.New("unknown command")
	}
}

func (a *App) RunScenario(id string, seconds float64, upgrades []model.UpgradeKind) (model.Summary, error) {
	if seconds <= 0 {
		seconds = 1
	}
	if _, err := a.StartRun(id); err != nil {
		return model.Summary{}, err
	}
	for _, kind := range upgrades {
		if _, err := a.ApplyUpgrade(id, kind); err != nil {
			return model.Summary{}, err
		}
	}
	if _, err := a.AdvanceRun(id, seconds); err != nil {
		return model.Summary{}, err
	}
	return a.Summary(id)
}

func (a *App) Preview(id string) (map[string]any, error) {
	run, err := a.store.GetRun(id)
	if err != nil {
		return nil, err
	}
	sim := engine.Restore(run)
	choices := upgrade.Choose(run.Player, run.Seed+int64(run.Revision), 3)
	return map[string]any{"run": run.ID, "difficulty": run.Wave, "choices": choices, "health": sim.State.Player.Health}, nil
}

func ValidateCommand(command Command) error {
	if command.RunID == "" {
		return errors.New("run id is required")
	}
	valid := map[string]bool{"start": true, "pause": true, "resume": true, "step": true, "upgrade": true, "stop": true, "summary": true}
	if !valid[command.Name] {
		return errors.New("unknown command")
	}
	if command.Name == "step" && command.Seconds <= 0 {
		return errors.New("step seconds must be positive")
	}
	if command.Name == "upgrade" {
		return upgrade.Validate(command.Upgrade)
	}
	return nil
}

func SortSummaries(summaries []model.Summary, ascending bool) []model.Summary {
	out := append([]model.Summary(nil), summaries...)
	sort.SliceStable(out, func(i, j int) bool {
		if ascending {
			return out[i].Score < out[j].Score
		}
		return out[i].Score > out[j].Score
	})
	return out
}

func MergeDispatchResults(results []model.DispatchResult) model.DispatchResult {
	out := model.DispatchResult{Records: make([]model.Record, 0), Summaries: make([]model.Summary, 0)}
	for _, result := range results {
		out.Records = append(out.Records, result.Records...)
		out.Summaries = append(out.Summaries, result.Summaries...)
		if out.Batch.ID == "" {
			out.Batch = result.Batch
		}
	}
	return out
}
