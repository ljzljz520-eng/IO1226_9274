package service

import (
	"miniarrow/internal/model"
	"path/filepath"
	"testing"
)

func TestWorkflowAccept(t *testing.T) {
	app, err := New(filepath.Join(t.TempDir(), "game.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()
	run, err := app.CreateRun("accept", 1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = app.StartRun(run.ID); err != nil {
		t.Fatal(err)
	}
	state, err := app.AdvanceRun(run.ID, 2)
	if err != nil {
		t.Fatal(err)
	}
	if state.Elapsed != 2 || state.Status != "running" {
		t.Fatal("accept workflow failed")
	}
}

func TestWorkflowPublish(t *testing.T) {
	app, err := New(filepath.Join(t.TempDir(), "game.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()
	a, _ := app.CreateRun("a", 1)
	b, _ := app.CreateRun("b", 2)
	result, err := app.DispatchMany([]string{a.ID, b.ID})
	if err != nil {
		t.Fatal(err)
	}
	if result.Batch.State != "complete" {
		t.Fatal("batch did not complete")
	}
}

func TestWorkflowReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "game.db")
	app, err := New(path)
	if err != nil {
		t.Fatal(err)
	}
	run, _ := app.CreateRun("reopen", 1)
	app.Close()
	app, err = New(path)
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()
	loaded, err := app.StartRun(run.ID)
	if err != nil || loaded.ID != run.ID {
		t.Fatal("reopen workflow failed")
	}
}

func TestUpgradeWorkflow(t *testing.T) {
	app, err := New(filepath.Join(t.TempDir(), "game.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()
	run, _ := app.CreateRun("upgrade", 4)
	app.StartRun(run.ID)
	state, err := app.ApplyUpgrade(run.ID, model.UpgradeBlast)
	if err != nil || len(state.Player.Upgrades) != 1 {
		t.Fatal("upgrade workflow failed")
	}
}
