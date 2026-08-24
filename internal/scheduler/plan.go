package scheduler

import (
	"miniarrow/internal/difficulty"
	"miniarrow/internal/model"
)

type Plan struct {
	RunID    string
	Steps    []string
	Priority int
	Estimate float64
}

func BuildPlan(state model.RunState) Plan {
	p := difficulty.Calculate(state.Elapsed, state.Player.Level, state.Wave)
	priority := state.Player.Level + state.Wave
	if state.Status == "ready" {
		priority += 2
	}
	return Plan{RunID: state.ID, Steps: []string{"load", "simulate", "score", "persist"}, Priority: priority, Estimate: difficulty.RecoveryWindow(p)}
}

func ValidatePlan(plan Plan) error {
	if plan.RunID == "" {
		return errPlan("run id is required")
	}
	if len(plan.Steps) < 4 {
		return errPlan("plan needs four steps")
	}
	if plan.Priority < 0 {
		return errPlan("priority cannot be negative")
	}
	return nil
}

type planError string

func (e planError) Error() string  { return string(e) }
func errPlan(message string) error { return planError(message) }

func Describe(plan Plan) string {
	return plan.RunID + ":" + plan.Steps[0] + ">" + plan.Steps[1] + ">" + plan.Steps[2] + ">" + plan.Steps[3]
}
