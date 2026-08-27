package service

import (
	"path/filepath"
	"testing"
)

func TestWorkflow27(t *testing.T) {
	app, err := New(filepath.Join(t.TempDir(), "game.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()
	run, err := app.CreateRun("double-dispatch", 27)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = app.StartRun(run.ID); err != nil {
		t.Fatal(err)
	}
	first, err := app.DispatchSample(run.ID)
	if err != nil {
		t.Fatal(err)
	}
	second, err := app.DispatchSample(run.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Summaries) != 1 || len(second.Summaries) != 1 {
		t.Fatalf("summary disappeared after consecutive dispatch: first=%d second=%d", len(first.Summaries), len(second.Summaries))
	}
	if first.Summaries[0].RunID != run.ID || second.Summaries[0].RunID != run.ID {
		t.Fatal("dispatch state leaked between samples")
	}
}
