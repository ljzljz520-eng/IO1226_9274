package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"miniarrow/internal/model"
	"miniarrow/internal/report"
	"miniarrow/internal/service"
	"os"
)

func main() {
	dbPath := flag.String("db", "miniarrow.db", "database path")
	name := flag.String("name", "现场弓手", "player name")
	seconds := flag.Float64("seconds", 10, "simulation seconds")
	flag.Parse()
	app, err := service.New(*dbPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer app.Close()
	run, err := app.CreateRun(*name, 42)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if _, err = app.StartRun(run.ID); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	state, err := app.AdvanceRun(run.ID, *seconds)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if _, err = app.ApplyUpgrade(run.ID, model.UpgradePierce); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	summary, err := app.Summary(run.ID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	data, _ := json.MarshalIndent(map[string]any{"state": state, "summary": summary, "line": report.Format(summary)}, "", "  ")
	fmt.Println(string(data))
}
