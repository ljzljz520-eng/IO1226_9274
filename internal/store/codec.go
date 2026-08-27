package store

import (
	"encoding/json"
	"miniarrow/internal/model"
)

func EncodeRun(run model.RunState) ([]byte, error) { return json.Marshal(run) }
func DecodeRun(data []byte) (model.RunState, error) {
	var run model.RunState
	err := json.Unmarshal(data, &run)
	return run, err
}
func EncodeRecord(record model.Record) ([]byte, error) { return json.Marshal(record) }
func DecodeRecord(data []byte) (model.Record, error) {
	var record model.Record
	err := json.Unmarshal(data, &record)
	return record, err
}
func EncodeBatch(batch model.Batch) ([]byte, error) { return json.Marshal(batch) }
func DecodeBatch(data []byte) (model.Batch, error) {
	var batch model.Batch
	err := json.Unmarshal(data, &batch)
	return batch, err
}
func EncodeAudit(audit model.Audit) ([]byte, error) { return json.Marshal(audit) }
func DecodeAudit(data []byte) (model.Audit, error) {
	var audit model.Audit
	err := json.Unmarshal(data, &audit)
	return audit, err
}
func EncodeProfile(profile model.Profile) ([]byte, error) { return json.Marshal(profile) }
func DecodeProfile(data []byte) (model.Profile, error) {
	var profile model.Profile
	err := json.Unmarshal(data, &profile)
	return profile, err
}
