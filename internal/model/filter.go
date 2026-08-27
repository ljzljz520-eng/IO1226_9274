package model

import "sort"

func SortEvents(events []Event) []Event {
	out := CloneEvents(events)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Tick == out[j].Tick {
			return out[i].Kind < out[j].Kind
		}
		return out[i].Tick < out[j].Tick
	})
	return out
}

func SortEnemies(enemies []Enemy) []Enemy {
	out := make([]Enemy, len(enemies))
	copy(out, enemies)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Health == out[j].Health {
			return out[i].ID < out[j].ID
		}
		return out[i].Health < out[j].Health
	})
	return out
}

func ActiveProjectiles(projectiles []Projectile) []Projectile {
	out := make([]Projectile, 0, len(projectiles))
	for _, projectile := range projectiles {
		if projectile.Age < 6 && projectile.Damage > 0 {
			out = append(out, CloneProjectile(projectile))
		}
	}
	return out
}

func EventWindow(events []Event, start, end int) []Event {
	if start < 0 {
		start = 0
	}
	if end > len(events) {
		end = len(events)
	}
	if start > end {
		start = end
	}
	return CloneEvents(events[start:end])
}
