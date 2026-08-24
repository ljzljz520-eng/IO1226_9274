package scheduler

import (
	"sort"
	"time"
)

type HistoryEntry struct {
	JobID    string
	RunID    string
	State    string
	Attempts int
	At       time.Time
	Message  string
}

type History struct{ entries []HistoryEntry }

func NewHistory() *History { return &History{entries: make([]HistoryEntry, 0)} }

func (h *History) Add(job Job, message string) {
	h.entries = append(h.entries, HistoryEntry{JobID: job.ID, RunID: job.RunID, State: job.State, Attempts: job.Attempts, At: time.Now().UTC(), Message: message})
}

func (h *History) AddState(job Job, state, message string) { job.State = state; h.Add(job, message) }

func (h *History) All() []HistoryEntry {
	out := make([]HistoryEntry, len(h.entries))
	copy(out, h.entries)
	return out
}

func (h *History) ForRun(runID string) []HistoryEntry {
	out := make([]HistoryEntry, 0)
	for _, entry := range h.entries {
		if entry.RunID == runID {
			out = append(out, entry)
		}
	}
	return out
}

func (h *History) Latest(runID string) (HistoryEntry, bool) {
	items := h.ForRun(runID)
	if len(items) == 0 {
		return HistoryEntry{}, false
	}
	return items[len(items)-1], true
}

func (h *History) Counts() map[string]int {
	counts := make(map[string]int)
	for _, entry := range h.entries {
		counts[entry.State]++
	}
	return counts
}

func (h *History) SortNewest() []HistoryEntry {
	out := h.All()
	sort.SliceStable(out, func(i, j int) bool { return out[i].At.After(out[j].At) })
	return out
}

func (h *History) Trim(limit int) {
	if limit < 0 {
		limit = 0
	}
	if len(h.entries) <= limit {
		return
	}
	h.entries = append([]HistoryEntry(nil), h.entries[len(h.entries)-limit:]...)
}

func (h *History) Size() int { return len(h.entries) }

func (h *History) HasState(runID, state string) bool {
	for _, entry := range h.entries {
		if entry.RunID == runID && entry.State == state {
			return true
		}
	}
	return false
}

func Backoff(attempts int) time.Duration {
	if attempts < 1 {
		attempts = 1
	}
	delay := time.Second * time.Duration(attempts*attempts)
	if delay > time.Minute {
		return time.Minute
	}
	return delay
}

func Retryable(state string, attempts int) bool {
	return (state == "failed" || state == "timeout") && attempts < 3
}
