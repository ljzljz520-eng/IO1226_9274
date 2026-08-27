package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"miniarrow/internal/model"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"go.etcd.io/bbolt"
)

var bucketNames = [][]byte{[]byte("runs"), []byte("records"), []byte("batches"), []byte("audits"), []byte("profiles")}

type Store struct {
	db   *bbolt.DB
	mu   sync.RWMutex
	path string
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("store path is required")
	}
	db, err := bbolt.Open(filepath.Clean(path), 0600, &bbolt.Options{Timeout: 2 * time.Second, NoSync: false})
	if err != nil {
		return nil, err
	}
	if err = db.Update(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if _, e := tx.CreateBucketIfNotExists(name); e != nil {
				return e
			}
		}
		return nil
	}); err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db, path: path}, nil
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) Path() string { return s.path }

func (s *Store) PutRun(run model.RunState) error { return s.put("runs", run.ID, run) }
func (s *Store) GetRun(id string) (model.RunState, error) {
	var run model.RunState
	err := s.get("runs", id, &run)
	return run, err
}
func (s *Store) DeleteRun(id string) error { return s.delete("runs", id) }
func (s *Store) ListRuns() ([]model.RunState, error) {
	var runs []model.RunState
	err := s.list("runs", &runs)
	sort.Slice(runs, func(i, j int) bool { return runs[i].UpdatedAt.Before(runs[j].UpdatedAt) })
	return runs, err
}

func (s *Store) PutRecord(record model.Record) error { return s.put("records", record.ID, record) }
func (s *Store) GetRecord(id string) (model.Record, error) {
	var record model.Record
	err := s.get("records", id, &record)
	return record, err
}
func (s *Store) ListRecords(runID string) ([]model.Record, error) {
	var records []model.Record
	if err := s.list("records", &records); err != nil {
		return nil, err
	}
	out := records[:0]
	for _, record := range records {
		if runID == "" || record.RunID == runID {
			out = append(out, record)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Sequence < out[j].Sequence })
	return out, nil
}

func (s *Store) PutBatch(batch model.Batch) error { return s.put("batches", batch.ID, batch) }
func (s *Store) GetBatch(id string) (model.Batch, error) {
	var batch model.Batch
	err := s.get("batches", id, &batch)
	return batch, err
}

func (s *Store) PutAudit(audit model.Audit) error { return s.put("audits", audit.ID, audit) }
func (s *Store) ListAudits(runID string) ([]model.Audit, error) {
	var audits []model.Audit
	if err := s.list("audits", &audits); err != nil {
		return nil, err
	}
	if runID == "" {
		return audits, nil
	}
	out := audits[:0]
	for _, audit := range audits {
		if audit.RunID == runID {
			out = append(out, audit)
		}
	}
	return out, nil
}

func (s *Store) PutProfile(profile model.Profile) error {
	return s.put("profiles", profile.ID, profile)
}
func (s *Store) GetProfile(id string) (model.Profile, error) {
	var profile model.Profile
	err := s.get("profiles", id, &profile)
	return profile, err
}

func (s *Store) put(bucket, key string, value any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s missing", bucket)
		}
		return b.Put([]byte(key), data)
	})
}

func (s *Store) get(bucket, key string, target any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s missing", bucket)
		}
		data := b.Get([]byte(key))
		if len(data) == 0 {
			return fmt.Errorf("%s %s not found", bucket, key)
		}
		return json.Unmarshal(data, target)
	})
}

func (s *Store) delete(bucket, key string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte(bucket)).Delete([]byte(key)) })
}

func (s *Store) list(bucket string, target any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	values := make([]json.RawMessage, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s missing", bucket)
		}
		return b.ForEach(func(_, value []byte) error { values = append(values, append([]byte(nil), value...)); return nil })
	})
	if err != nil {
		return err
	}
	encoded, err := json.Marshal(values)
	if err != nil {
		return err
	}
	return json.Unmarshal(encoded, target)
}

func (s *Store) Count(bucket string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return 0, errors.New("store closed")
	}
	count := 0
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s missing", bucket)
		}
		count = b.Stats().KeyN
		return nil
	})
	return count, err
}
