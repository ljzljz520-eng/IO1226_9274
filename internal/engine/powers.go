package engine

import (
	"math"
	"miniarrow/internal/model"
)

func (s *Simulation) piercePower() int {
	if !model.HasUpgrade(s.State.Player, model.UpgradePierce) {
		return 0
	}
	return 1 + s.State.Player.Level/6
}

func (s *Simulation) bouncePower() int {
	if !model.HasUpgrade(s.State.Player, model.UpgradeBounce) {
		return 0
	}
	return 1 + s.State.Player.Level/8
}

func (s *Simulation) blastPower() float64 {
	if !model.HasUpgrade(s.State.Player, model.UpgradeBlast) {
		return 0
	}
	return 5 + float64(s.State.Player.Level)*0.4
}

func (s *Simulation) fireCooldown() float64 {
	cooldown := 1.2 - float64(s.State.Player.Level)*0.025
	if model.HasUpgrade(s.State.Player, model.UpgradeOrbit) {
		cooldown -= 0.1
	}
	if cooldown < 0.35 {
		return 0.35
	}
	return cooldown
}

func (s *Simulation) OrbitDamage(angle float64) float64 {
	if !model.HasUpgrade(s.State.Player, model.UpgradeOrbit) {
		return 0
	}
	return 5 + math.Abs(math.Sin(angle))*float64(s.State.Player.Level)
}

func (s *Simulation) ApplyUpgrade(kind model.UpgradeKind) bool {
	if s.State.Status == "finished" || model.HasUpgrade(s.State.Player, kind) {
		return false
	}
	model.AddUpgrade(&s.State.Player, kind)
	s.emit("upgrade", s.State.Player.ID, float64(model.UpgradeCount(s.State.Player)), string(kind))
	return true
}

func (s *Simulation) UpgradeChoices() []model.UpgradeKind {
	choices := make([]model.UpgradeKind, 0, 4)
	for _, kind := range []model.UpgradeKind{model.UpgradePierce, model.UpgradeBounce, model.UpgradeBlast, model.UpgradeOrbit} {
		if !model.HasUpgrade(s.State.Player, kind) {
			choices = append(choices, kind)
		}
	}
	return choices
}
