package upgrade

import "miniarrow/internal/model"

type Progress struct {
	Current   int
	Next      int
	Remaining int
	Percent   float64
}

func Calculate(experience, level int) Progress {
	if level < 1 {
		level = 1
	}
	base := 80 + (level-1)*45 + (level-1)*(level-1)*10
	previous := 0
	if level > 1 {
		previous = 80 + (level-2)*45 + (level-2)*(level-2)*10
	}
	span := base - previous
	current := experience - previous
	if current < 0 {
		current = 0
	}
	if current > span {
		current = span
	}
	percent := float64(current) / float64(span) * 100
	return Progress{Current: current, Next: base, Remaining: span - current, Percent: percent}
}

func CanSelect(player model.Player, kind model.UpgradeKind) bool {
	if model.HasUpgrade(player, kind) {
		return false
	}
	_, ok := Find(kind)
	return ok && player.Level >= 1
}

func Describe(kind model.UpgradeKind) string {
	if item, ok := Find(kind); ok {
		return item.Name + ": " + item.Description
	}
	return "unknown upgrade"
}
