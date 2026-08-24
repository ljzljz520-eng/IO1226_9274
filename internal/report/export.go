package report

import (
	"bytes"
	"encoding/csv"
	"miniarrow/internal/model"
	"strconv"
)

func CSV(summaries []model.Summary) ([]byte, error) {
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)
	if err := writer.Write([]string{"run_id", "name", "elapsed", "level", "wave", "score", "kills", "difficulty"}); err != nil {
		return nil, err
	}
	for _, summary := range summaries {
		if err := writer.Write([]string{summary.RunID, summary.DisplayName, strconv.FormatFloat(summary.Elapsed, 'f', 1, 64), strconv.Itoa(summary.Level), strconv.Itoa(summary.Wave), strconv.Itoa(summary.Score), strconv.Itoa(summary.Kills), summary.Difficulty}); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	return buffer.Bytes(), writer.Error()
}

func EventCounts(events []model.Event) map[string]int {
	counts := make(map[string]int)
	for _, event := range events {
		counts[event.Kind]++
	}
	return counts
}

func FilterEvents(events []model.Event, kind string) []model.Event {
	out := make([]model.Event, 0)
	for _, event := range events {
		if kind == "" || event.Kind == kind {
			out = append(out, event)
		}
	}
	return out
}
