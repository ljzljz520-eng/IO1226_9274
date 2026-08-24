package report

import (
	"encoding/json"
	"fmt"
	"miniarrow/internal/difficulty"
	"miniarrow/internal/engine"
	"miniarrow/internal/model"
	"sort"
	"strings"
)

func BuildSummary(run model.RunState, events []model.Event) model.Summary {
	metrics := engine.Measure(events)
	p := difficulty.Calculate(run.Elapsed, run.Player.Level, run.Wave)
	latest := events
	if len(latest) > 6 {
		latest = latest[len(latest)-6:]
	}
	return model.Summary{RunID: run.ID, DisplayName: run.Player.Name, Elapsed: run.Elapsed, Level: run.Player.Level, Wave: run.Wave, Score: run.Player.Score, Kills: metrics.Kills, Status: run.Status, Upgrades: append([]model.UpgradeKind(nil), run.Player.Upgrades...), Difficulty: difficulty.Label(p), LatestEvents: model.CloneEvents(latest)}
}

func Format(summary model.Summary) string {
	return fmt.Sprintf("%s | %s | %.1fs | lvl %d | wave %d | score %d | kills %d | %s | upgrades=%s", summary.RunID, summary.DisplayName, summary.Elapsed, summary.Level, summary.Wave, summary.Score, summary.Kills, summary.Difficulty, strings.Join(upgradeNames(summary.Upgrades), ","))
}

func upgradeNames(upgrades []model.UpgradeKind) []string {
	out := make([]string, len(upgrades))
	for i, upgrade := range upgrades {
		out[i] = string(upgrade)
	}
	return out
}

func Marshal(summary model.Summary) ([]byte, error) { return json.MarshalIndent(summary, "", "  ") }

func Rank(summaries []model.Summary) []model.Summary {
	out := append([]model.Summary(nil), summaries...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Score == out[j].Score {
			return out[i].Elapsed > out[j].Elapsed
		}
		return out[i].Score > out[j].Score
	})
	return out
}

func Leaderboard(summaries []model.Summary, limit int) []model.LeaderboardEntry {
	ranked := Rank(summaries)
	if limit <= 0 || limit > len(ranked) {
		limit = len(ranked)
	}
	out := make([]model.LeaderboardEntry, 0, limit)
	for _, summary := range ranked[:limit] {
		out = append(out, model.LeaderboardEntry{RunID: summary.RunID, Name: summary.DisplayName, Score: summary.Score, Level: summary.Level, Elapsed: summary.Elapsed})
	}
	return out
}

func Compare(left, right model.Summary) int {
	if left.Score != right.Score {
		if left.Score > right.Score {
			return 1
		}
		return -1
	}
	if left.Elapsed > right.Elapsed {
		return 1
	}
	if left.Elapsed < right.Elapsed {
		return -1
	}
	return 0
}

func StatusLine(run model.RunState) string {
	switch run.Status {
	case "running":
		return "进行中"
	case "paused":
		return "已暂停"
	case "finished":
		return "已结束"
	default:
		return "待开始"
	}
}
