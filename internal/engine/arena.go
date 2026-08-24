package engine

import (
	"fmt"
	"math"
	"miniarrow/internal/difficulty"
	"miniarrow/internal/model"
	"sort"
)

type Arena struct {
	Width          float64
	Height         float64
	Center         model.Vec2
	SpawnRadius    float64
	BoundaryDamage float64
}

func DefaultArena() Arena {
	return Arena{Width: 80, Height: 80, Center: model.Vec2{}, SpawnRadius: 25, BoundaryDamage: 1}
}

func NewArena(width, height float64) Arena {
	if width < 20 {
		width = 20
	}
	if height < 20 {
		height = 20
	}
	return Arena{Width: width, Height: height, Center: model.Vec2{}, SpawnRadius: math.Min(width, height) * 0.4, BoundaryDamage: 1}
}

func (a Arena) Contains(position model.Vec2) bool {
	return math.Abs(position.X-a.Center.X) <= a.Width/2 && math.Abs(position.Y-a.Center.Y) <= a.Height/2
}

func (a Arena) Clamp(position model.Vec2) model.Vec2 {
	if position.X < -a.Width/2 {
		position.X = -a.Width / 2
	}
	if position.X > a.Width/2 {
		position.X = a.Width / 2
	}
	if position.Y < -a.Height/2 {
		position.Y = -a.Height / 2
	}
	if position.Y > a.Height/2 {
		position.Y = a.Height / 2
	}
	return position
}

func (a Arena) SpawnPosition(index int, seed int64) model.Vec2 {
	angle := float64((int(seed)+index*53)%360) * math.Pi / 180
	distance := a.SpawnRadius + float64(index%5)
	return a.Center.Add(model.Vec2{X: math.Cos(angle) * distance, Y: math.Sin(angle) * distance})
}

func (a Arena) DistanceFromCenter(position model.Vec2) float64 {
	return position.Sub(a.Center).Length()
}

func (a Arena) EdgePenalty(position model.Vec2) float64 {
	distance := a.DistanceFromCenter(position)
	limit := math.Min(a.Width, a.Height) / 2
	if distance <= limit*0.7 {
		return 0
	}
	penalty := (distance - limit*0.7) / (limit * 0.3)
	if penalty > 1 {
		penalty = 1
	}
	return penalty
}

func MovePlayer(state *model.RunState, arena Arena, direction model.Vec2, dt float64) {
	if dt <= 0 {
		return
	}
	speed := 6.0
	if model.HasUpgrade(state.Player, model.UpgradeOrbit) {
		speed += 0.8
	}
	state.Player.Position = arena.Clamp(state.Player.Position.Add(direction.Normalize().Scale(speed * dt)))
}

func ApplyBoundaryPressure(state *model.RunState, arena Arena, dt float64) float64 {
	penalty := arena.EdgePenalty(state.Player.Position)
	if penalty <= 0 || dt <= 0 {
		return 0
	}
	damage := penalty * arena.BoundaryDamage * dt
	state.Player.Health -= damage
	return damage
}

func Formation(enemies []model.Enemy) map[model.EnemyKind]int {
	counts := make(map[model.EnemyKind]int)
	for _, enemy := range enemies {
		if enemy.Health > 0 {
			counts[enemy.Kind]++
		}
	}
	return counts
}

func ThreatOrder(enemies []model.Enemy) []model.Enemy {
	out := make([]model.Enemy, len(enemies))
	copy(out, enemies)
	sort.SliceStable(out, func(i, j int) bool {
		left := out[i].Attack * out[i].Speed
		right := out[j].Attack * out[j].Speed
		if left == right {
			return out[i].ID < out[j].ID
		}
		return left > right
	})
	return out
}

func BuildWave(wave int, elapsed float64, level int, arena Arena, seed int64) []model.Enemy {
	profile := difficulty.Calculate(elapsed, level, wave)
	count := difficulty.SpawnBudget(profile)
	out := make([]model.Enemy, 0, count)
	for i := 0; i < count; i++ {
		kind := model.EnemyScout
		if difficulty.IsBossWave(wave) && i == 0 {
			kind = model.EnemyCaptain
		} else if i%4 == 0 {
			kind = model.EnemyBrute
		} else if i%3 == 0 {
			kind = model.EnemyWisp
		}
		health, speed, attack := difficulty.EnemyStats(string(kind), profile)
		out = append(out, model.Enemy{ID: fmt.Sprintf("wave-%d-%d", wave, i), Kind: kind, Position: arena.SpawnPosition(i, seed), Health: health, MaxHealth: health, Speed: speed, Attack: attack})
	}
	return out
}

func SimulateWave(sim *Simulation, arena Arena, seconds, step float64) []model.Event {
	if step <= 0 {
		step = 0.5
	}
	if seconds < 0 {
		seconds = 0
	}
	if sim.State.Status == "ready" {
		sim.Start()
	}
	for elapsed := 0.0; elapsed < seconds && sim.State.Status == "running"; elapsed += step {
		ApplyBoundaryPressure(&sim.State, arena, step)
		sim.Step(step)
	}
	return model.CloneEvents(sim.Events)
}

func ArenaReport(state model.RunState, arena Arena) map[string]any {
	return map[string]any{"inside": arena.Contains(state.Player.Position), "edge_penalty": arena.EdgePenalty(state.Player.Position), "enemies": Formation(state.Enemies), "threat_order": ThreatOrder(state.Enemies)}
}
