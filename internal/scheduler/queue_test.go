package scheduler

import "testing"

func TestQueuePriority(t *testing.T) {
	q := NewQueue()
	q.Enqueue("low", 1)
	q.Enqueue("high", 5)
	jobs := q.Claim(2)
	if len(jobs) != 2 || jobs[0].RunID != "high" {
		t.Fatal("priority ordering failed")
	}
}

func TestPlanValidation(t *testing.T) {
	plan := Plan{RunID: "r", Steps: []string{"a", "b", "c", "d"}, Priority: 1}
	if err := ValidatePlan(plan); err != nil {
		t.Fatal(err)
	}
}
