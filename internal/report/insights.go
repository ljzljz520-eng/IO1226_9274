package report

import (
	"math"
	"miniarrow/internal/engine"
	"miniarrow/internal/model"
	"sort"
)

type Insight struct {
	Code     string
	Title    string
	Detail   string
	Severity string
	Value    float64
}

func BuildInsights(run model.RunState, events []model.Event) []Insight {
	metrics := engine.Measure(events)
	out := make([]Insight, 0, 8)
	if metrics.Kills == 0 && run.Elapsed > 30 {
		out = append(out, Insight{Code: "no-kills", Title: "没有击杀", Detail: "生存期间没有完成击杀", Severity: "warning", Value: run.Elapsed})
	}
	if metrics.DamageTaken > metrics.DamageDealt {
		out = append(out, Insight{Code: "pressure", Title: "承伤偏高", Detail: "受到的伤害超过造成的伤害", Severity: "danger", Value: metrics.DamageTaken - metrics.DamageDealt})
	}
	if model.HasUpgrade(run.Player, model.UpgradePierce) {
		out = append(out, Insight{Code: "pierce", Title: "穿透箭", Detail: "箭矢可穿透目标", Severity: "info", Value: float64(run.Player.Level)})
	}
	if model.HasUpgrade(run.Player, model.UpgradeBounce) {
		out = append(out, Insight{Code: "bounce", Title: "弹跳箭", Detail: "箭矢会在命中后改变方向", Severity: "info", Value: float64(run.Player.Level)})
	}
	if model.HasUpgrade(run.Player, model.UpgradeBlast) {
		out = append(out, Insight{Code: "blast", Title: "范围爆炸", Detail: "命中会影响附近敌人", Severity: "info", Value: float64(run.Player.Level)})
	}
	if model.HasUpgrade(run.Player, model.UpgradeOrbit) {
		out = append(out, Insight{Code: "orbit", Title: "旋转箭阵", Detail: "旋转箭阵提升持续输出", Severity: "info", Value: float64(run.Player.Level)})
	}
	if run.Status == "finished" && run.Player.Health <= 0 {
		out = append(out, Insight{Code: "defeat", Title: "防线失守", Detail: "角色生命值归零", Severity: "danger", Value: run.Elapsed})
	}
	return out
}

func Trend(summaries []model.Summary) string {
	if len(summaries) < 2 {
		return "insufficient"
	}
	ordered := Rank(summaries)
	if ordered[0].Score == ordered[len(ordered)-1].Score {
		return "flat"
	}
	if ordered[0].Elapsed >= ordered[len(ordered)-1].Elapsed {
		return "improving"
	}
	return "mixed"
}

func Percentile(values []float64, percentile float64) float64 {
	if len(values) == 0 {
		return 0
	}
	items := append([]float64(nil), values...)
	sort.Float64s(items)
	if percentile < 0 {
		percentile = 0
	}
	if percentile > 1 {
		percentile = 1
	}
	index := int(math.Round(percentile * float64(len(items)-1)))
	return items[index]
}

func ScoreBand(score int) string {
	if score >= 1000 {
		return "S"
	}
	if score >= 500 {
		return "A"
	}
	if score >= 200 {
		return "B"
	}
	if score >= 80 {
		return "C"
	}
	return "D"
}

func Aggregate(summaries []model.Summary) map[string]float64 {
	result := map[string]float64{"runs": float64(len(summaries))}
	if len(summaries) == 0 {
		return result
	}
	var score, elapsed, kills float64
	for _, summary := range summaries {
		score += float64(summary.Score)
		elapsed += summary.Elapsed
		kills += float64(summary.Kills)
	}
	result["score_avg"] = score / float64(len(summaries))
	result["elapsed_avg"] = elapsed / float64(len(summaries))
	result["kills_avg"] = kills / float64(len(summaries))
	return result
}

func GroupByDifficulty(summaries []model.Summary) map[string][]model.Summary {
	groups := make(map[string][]model.Summary)
	for _, summary := range summaries {
		groups[summary.Difficulty] = append(groups[summary.Difficulty], summary)
	}
	return groups
}
