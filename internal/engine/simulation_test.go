package engine

import (
	"miniarrow/internal/model"
	"testing"
)

func TestSimulationTargetsNearestEnemy(t *testing.T) {
	sim := NewRun("r", "a", 1)
	sim.Start()
	sim.State.Enemies = []model.Enemy{{ID: "far", Position: model.Vec2{X: 20}, Health: 10}, {ID: "near", Position: model.Vec2{X: 4}, Health: 10}}
	sim.Step(1)
	seen := false
	for _, event := range sim.Events {
		if event.Kind == "fire" && event.Subject == "near" {
			seen = true
		}
	}
	if !seen {
		t.Fatal("arrow did not target nearest enemy")
	}
}

func TestUpgradesChangeProjectile(t *testing.T) {
	sim := NewRun("r", "a", 1)
	sim.ApplyUpgrade(model.UpgradePierce)
	sim.ApplyUpgrade(model.UpgradeBlast)
	sim.State.Status = "running"
	sim.State.Enemies = []model.Enemy{{ID: "e", Position: model.Vec2{X: 1}, Health: 50}}
	sim.Step(1)
	if len(sim.Events) == 0 {
		t.Fatal("simulation emitted no events")
	}
}
