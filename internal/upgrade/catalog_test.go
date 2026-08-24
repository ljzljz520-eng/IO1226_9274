package upgrade

import (
	"miniarrow/internal/model"
	"testing"
)

func TestCatalogHasFourPowers(t *testing.T) {
	if len(All()) != 4 {
		t.Fatal("expected four upgrades")
	}
}

func TestChooseExcludesOwned(t *testing.T) {
	p := model.Player{Upgrades: []model.UpgradeKind{model.UpgradePierce}}
	choices := Choose(p, 3, 4)
	for _, choice := range choices {
		if choice.Kind == model.UpgradePierce {
			t.Fatal("owned upgrade returned")
		}
	}
}
