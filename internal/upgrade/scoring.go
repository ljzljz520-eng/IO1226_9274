package upgrade

import (
	"miniarrow/internal/model"
	"sort"
)

type ChoiceScore struct {
	Definition Definition
	Score      float64
	Reason     string
}

func ScoreChoices(player model.Player, elapsed float64) []ChoiceScore {
	choices := Eligible(player)
	out := make([]ChoiceScore, 0, len(choices))
	for _, choice := range choices {
		score := float64(choice.Weight) + float64(player.Level)*0.2
		reason := "适合当前生存阶段"
		switch choice.Kind {
		case model.UpgradePierce:
			score += elapsed / 90
			reason = "敌人密度提升时收益稳定"
		case model.UpgradeBounce:
			score += float64(player.Level) / 4
			reason = "单发箭矢可覆盖多条路线"
		case model.UpgradeBlast:
			score += elapsed / 70
			reason = "范围伤害缓解包围"
		case model.UpgradeOrbit:
			score += 1
			reason = "降低冷却并保护近身"
		}
		out = append(out, ChoiceScore{Definition: choice, Score: score, Reason: reason})
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out
}

func BestChoice(player model.Player, elapsed float64) (ChoiceScore, bool) {
	choices := ScoreChoices(player, elapsed)
	if len(choices) == 0 {
		return ChoiceScore{}, false
	}
	return choices[0], true
}

func UpgradePower(player model.Player, kind model.UpgradeKind) float64 {
	if !model.HasUpgrade(player, kind) {
		return 0
	}
	level := float64(player.Level)
	switch kind {
	case model.UpgradePierce:
		return 1 + level/6
	case model.UpgradeBounce:
		return 1 + level/8
	case model.UpgradeBlast:
		return 5 + level*0.4
	case model.UpgradeOrbit:
		return 5 + level
	}
	return 0
}

func UpgradeSetPower(player model.Player) float64 {
	total := 0.0
	for _, kind := range player.Upgrades {
		total += UpgradePower(player, kind)
	}
	return total
}

func MissingKinds(player model.Player) []model.UpgradeKind {
	out := make([]model.UpgradeKind, 0)
	for _, item := range catalog {
		if !model.HasUpgrade(player, item.Kind) {
			out = append(out, item.Kind)
		}
	}
	return out
}

func Complete(player model.Player) bool { return len(MissingKinds(player)) == 0 }
