package model

import "testing"

func TestCloneRunIsIndependent(t *testing.T) {
	run := RunState{Player: Player{Upgrades: []UpgradeKind{UpgradePierce}}, Enemies: []Enemy{{ID: "e"}}}
	copy := CloneRun(run)
	copy.Player.Upgrades[0] = UpgradeBounce
	copy.Enemies[0].ID = "changed"
	if run.Player.Upgrades[0] != UpgradePierce || run.Enemies[0].ID != "e" {
		t.Fatal("clone shares mutable data")
	}
}

func TestVectorDistance(t *testing.T) {
	if (Vec2{X: 3}).Distance(Vec2{}) != 9 {
		t.Fatal("distance should be squared")
	}
}

func TestUpgradeHelpers(t *testing.T) {
	p := Player{}
	AddUpgrade(&p, UpgradeBlast)
	AddUpgrade(&p, UpgradeBlast)
	if UpgradeCount(p) != 1 || !HasUpgrade(p, UpgradeBlast) {
		t.Fatal("upgrade helper failed")
	}
}
