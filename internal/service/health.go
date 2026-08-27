package service

import "miniarrow/internal/model"

type Health struct {
	Runs     int
	Records  int
	Batches  int
	Audits   int
	Profiles int
}

func (a *App) Health() (Health, error) {
	runs, err := a.store.Count("runs")
	if err != nil {
		return Health{}, err
	}
	records, err := a.store.Count("records")
	if err != nil {
		return Health{}, err
	}
	batches, err := a.store.Count("batches")
	if err != nil {
		return Health{}, err
	}
	audits, err := a.store.Count("audits")
	if err != nil {
		return Health{}, err
	}
	profiles, err := a.store.Count("profiles")
	if err != nil {
		return Health{}, err
	}
	return Health{Runs: runs, Records: records, Batches: batches, Audits: audits, Profiles: profiles}, nil
}

func (a *App) LatestRecord(id string) (model.Record, error) {
	records, err := a.store.ListRecords(id)
	if err != nil {
		return model.Record{}, err
	}
	if len(records) == 0 {
		return model.Record{}, errMissing("record")
	}
	return records[len(records)-1], nil
}

type missingError string

func (e missingError) Error() string { return string(e) + " missing" }
func errMissing(value string) error  { return missingError(value) }
