package report

import (
	"miniarrow/internal/model"
	"testing"
)

func TestSummaryAndRanking(t *testing.T) {
	run := model.RunState{ID: "r", Player: model.Player{Name: "a", Level: 2, Score: 40, Upgrades: []model.UpgradeKind{model.UpgradeBounce}}, Elapsed: 120, Wave: 4, Status: "running"}
	summary := BuildSummary(run, []model.Event{{Kind: "kill"}, {Kind: "kill"}})
	if summary.Kills != 2 || summary.Score != 40 {
		t.Fatal("summary incorrect")
	}
	ranked := Rank([]model.Summary{summary, {Score: 1}})
	if ranked[0].Score != 40 {
		t.Fatal("ranking incorrect")
	}
}

func TestCSVExport(t *testing.T) {
	data, err := CSV([]model.Summary{{RunID: "r", DisplayName: "a"}})
	if err != nil || len(data) == 0 {
		t.Fatal("csv export failed")
	}
}
