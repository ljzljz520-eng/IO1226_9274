package store

import (
	"miniarrow/internal/model"
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "runs.db")
	first, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	run := model.RunState{ID: "persist", Status: "running", Player: model.Player{Name: "saved"}}
	if err = first.PutRun(run); err != nil {
		t.Fatal(err)
	}
	if err = first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	loaded, err := second.GetRun("persist")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Player.Name != "saved" {
		t.Fatalf("unexpected persisted name: %s", loaded.Player.Name)
	}
}

func TestStoreListsRecordsByRun(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "runs.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for i := 0; i < 3; i++ {
		if err := db.PutRecord(model.Record{ID: string(rune('a' + i)), RunID: "r", Sequence: i}); err != nil {
			t.Fatal(err)
		}
	}
	items, err := db.ListRecords("r")
	if err != nil || len(items) != 3 {
		t.Fatalf("records missing: %v %d", err, len(items))
	}
}
