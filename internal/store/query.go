package store

import (
	"miniarrow/internal/model"
	"sort"
)

type RunQuery struct {
	Status   string
	MinLevel int
	MaxLevel int
	MinScore int
	Name     string
	Limit    int
}

func (s *Store) QueryRuns(query RunQuery) ([]model.RunState, error) {
	runs, err := s.ListRuns()
	if err != nil {
		return nil, err
	}
	out := make([]model.RunState, 0, len(runs))
	for _, run := range runs {
		if query.Status != "" && run.Status != query.Status {
			continue
		}
		if query.MinLevel > 0 && run.Player.Level < query.MinLevel {
			continue
		}
		if query.MaxLevel > 0 && run.Player.Level > query.MaxLevel {
			continue
		}
		if query.MinScore > 0 && run.Player.Score < query.MinScore {
			continue
		}
		if query.Name != "" && run.Player.Name != query.Name {
			continue
		}
		out = append(out, run)
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Player.Score > out[j].Player.Score })
	if query.Limit > 0 && len(out) > query.Limit {
		out = out[:query.Limit]
	}
	return out, nil
}

func (s *Store) LatestRuns(limit int) ([]model.RunState, error) {
	return s.QueryRuns(RunQuery{Limit: limit})
}

func (s *Store) RecordsForBatch(batchID string) ([]model.Record, error) {
	batch, err := s.GetBatch(batchID)
	if err != nil {
		return nil, err
	}
	all, err := s.ListRecords("")
	if err != nil {
		return nil, err
	}
	allowed := make(map[string]bool)
	for _, id := range batch.RunIDs {
		allowed[id] = true
	}
	out := make([]model.Record, 0)
	for _, record := range all {
		if allowed[record.RunID] && record.Sequence == batch.Sequence {
			out = append(out, record)
		}
	}
	return out, nil
}

func (s *Store) AuditTrail(runID string) ([]model.Audit, error) {
	audits, err := s.ListAudits(runID)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(audits, func(i, j int) bool { return audits[i].At.Before(audits[j].At) })
	return audits, nil
}

func (s *Store) DeleteRecords(runID string) error {
	records, err := s.ListRecords(runID)
	if err != nil {
		return err
	}
	for _, record := range records {
		if err := s.delete("records", record.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) SnapshotCounts() (map[string]int, error) {
	counts := make(map[string]int)
	for _, bucket := range []string{"runs", "records", "batches", "audits", "profiles"} {
		value, err := s.Count(bucket)
		if err != nil {
			return nil, err
		}
		counts[bucket] = value
	}
	return counts, nil
}
