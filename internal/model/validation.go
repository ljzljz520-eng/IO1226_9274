package model

import (
	"errors"
	"fmt"
)

func ValidateRun(run RunState) error {
	if run.ID == "" {
		return errors.New("run id is required")
	}
	if run.Player.Name == "" {
		return errors.New("player name is required")
	}
	if run.Player.MaxHealth < 0 || run.Player.Health > run.Player.MaxHealth {
		return errors.New("invalid health range")
	}
	if run.Player.Level < 1 {
		return errors.New("level must be positive")
	}
	if run.Elapsed < 0 {
		return errors.New("elapsed time cannot be negative")
	}
	if run.Wave < 1 {
		return errors.New("wave must be positive")
	}
	seen := make(map[string]bool)
	for _, enemy := range run.Enemies {
		if enemy.ID == "" {
			return errors.New("enemy id is required")
		}
		if seen[enemy.ID] {
			return fmt.Errorf("duplicate enemy %s", enemy.ID)
		}
		seen[enemy.ID] = true
		if enemy.MaxHealth < 0 || enemy.Health > enemy.MaxHealth {
			return fmt.Errorf("invalid enemy health %s", enemy.ID)
		}
	}
	for _, projectile := range run.Projectiles {
		if projectile.ID == "" {
			return errors.New("projectile id is required")
		}
		if projectile.Damage < 0 {
			return errors.New("projectile damage cannot be negative")
		}
	}
	return nil
}

func ValidateRecord(record Record) error {
	if record.ID == "" || record.RunID == "" {
		return errors.New("record identifiers are required")
	}
	if record.Sequence < 1 {
		return errors.New("record sequence must be positive")
	}
	if record.Elapsed < 0 {
		return errors.New("record elapsed cannot be negative")
	}
	return nil
}

func ValidateBatch(batch Batch) error {
	if batch.ID == "" {
		return errors.New("batch id is required")
	}
	if len(batch.RunIDs) == 0 {
		return errors.New("batch needs runs")
	}
	for _, id := range batch.RunIDs {
		if id == "" {
			return errors.New("batch has empty run id")
		}
	}
	return nil
}

func ValidateAudit(audit Audit) error {
	if audit.ID == "" || audit.Action == "" {
		return errors.New("audit id and action are required")
	}
	if audit.RunID == "" && audit.BatchID == "" {
		return errors.New("audit needs a subject")
	}
	return nil
}

func NormalizeRun(run RunState) RunState {
	out := CloneRun(run)
	if out.Player.MaxHealth == 0 {
		out.Player.MaxHealth = 100
	}
	if out.Player.Health == 0 && out.Status == "ready" {
		out.Player.Health = out.Player.MaxHealth
	}
	if out.Player.Level < 1 {
		out.Player.Level = 1
	}
	if out.Wave < 1 {
		out.Wave = 1
	}
	if out.Status == "" {
		out.Status = "ready"
	}
	return out
}

func EventValue(events []Event, kind string) float64 {
	var value float64
	for _, event := range events {
		if event.Kind == kind {
			value += event.Value
		}
	}
	return value
}

func EventSubjects(events []Event, kind string) []string {
	out := make([]string, 0)
	seen := make(map[string]bool)
	for _, event := range events {
		if event.Kind == kind && !seen[event.Subject] {
			out = append(out, event.Subject)
			seen[event.Subject] = true
		}
	}
	return out
}

func MergeEvents(groups ...[]Event) []Event {
	total := 0
	for _, group := range groups {
		total += len(group)
	}
	out := make([]Event, 0, total)
	for _, group := range groups {
		out = append(out, group...)
	}
	return out
}
