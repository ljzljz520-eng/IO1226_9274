package engine

import "miniarrow/internal/model"

type Metrics struct {
	Kills       int
	DamageDealt float64
	DamageTaken float64
	Shots       int
	Upgrades    int
}

func Measure(events []model.Event) Metrics {
	var out Metrics
	for _, event := range events {
		switch event.Kind {
		case "kill":
			out.Kills++
		case "hit", "blast":
			out.DamageDealt += event.Value
		case "damage":
			out.DamageTaken += event.Value
		case "fire":
			out.Shots++
		case "upgrade":
			out.Upgrades++
		}
	}
	return out
}

func MergeMetrics(groups ...Metrics) Metrics {
	var out Metrics
	for _, group := range groups {
		out.Kills += group.Kills
		out.DamageDealt += group.DamageDealt
		out.DamageTaken += group.DamageTaken
		out.Shots += group.Shots
		out.Upgrades += group.Upgrades
	}
	return out
}

func SurvivalRank(metrics Metrics, elapsed float64) string {
	if metrics.Kills > 80 && elapsed > 300 {
		return "legend"
	}
	if metrics.Kills > 35 && elapsed > 180 {
		return "veteran"
	}
	if metrics.Kills > 10 || elapsed > 90 {
		return "survivor"
	}
	return "rookie"
}
