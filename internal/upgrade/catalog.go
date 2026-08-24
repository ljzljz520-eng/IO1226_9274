package upgrade

import (
	"fmt"
	"miniarrow/internal/model"
)

type Definition struct {
	Kind        model.UpgradeKind
	Name        string
	Description string
	Weight      int
	MaxLevel    int
}

var catalog = []Definition{
	{Kind: model.UpgradePierce, Name: "穿透箭", Description: "箭矢穿过更多敌人", Weight: 4, MaxLevel: 5},
	{Kind: model.UpgradeBounce, Name: "弹跳箭", Description: "箭矢命中后寻找下一个目标", Weight: 3, MaxLevel: 5},
	{Kind: model.UpgradeBlast, Name: "范围爆炸", Description: "命中造成范围伤害", Weight: 2, MaxLevel: 4},
	{Kind: model.UpgradeOrbit, Name: "旋转箭阵", Description: "旋转箭阵降低冷却并反击", Weight: 1, MaxLevel: 3},
}

func All() []Definition { out := make([]Definition, len(catalog)); copy(out, catalog); return out }

func Find(kind model.UpgradeKind) (Definition, bool) {
	for _, item := range catalog {
		if item.Kind == kind {
			return item, true
		}
	}
	return Definition{}, false
}

func Validate(kind model.UpgradeKind) error {
	if _, ok := Find(kind); !ok {
		return fmt.Errorf("unknown upgrade %q", kind)
	}
	return nil
}

func Eligible(player model.Player) []Definition {
	out := make([]Definition, 0, len(catalog))
	for _, item := range catalog {
		if !model.HasUpgrade(player, item.Kind) {
			out = append(out, item)
		}
	}
	return out
}

func Choose(player model.Player, seed int64, count int) []Definition {
	available := Eligible(player)
	if count <= 0 || len(available) == 0 {
		return nil
	}
	if count > len(available) {
		count = len(available)
	}
	start := int(seed % int64(len(available)))
	if start < 0 {
		start += len(available)
	}
	out := make([]Definition, 0, count)
	for i := 0; i < count; i++ {
		out = append(out, available[(start+i)%len(available)])
	}
	return out
}

func Summary(player model.Player) []string {
	out := make([]string, 0, len(player.Upgrades))
	for _, kind := range player.Upgrades {
		if item, ok := Find(kind); ok {
			out = append(out, item.Name)
		}
	}
	return out
}
