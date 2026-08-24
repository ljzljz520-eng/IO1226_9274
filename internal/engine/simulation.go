package engine

import (
	"fmt"
	"math"
	"miniarrow/internal/difficulty"
	"miniarrow/internal/model"
	"time"
)

type Simulation struct {
	State  model.RunState
	Events []model.Event
	tick   int
}

func NewRun(id, name string, seed int64) *Simulation {
	return &Simulation{State: model.RunState{ID: id, Player: model.Player{ID: "player-" + id, Name: name, Position: model.Vec2{}, Health: 100, MaxHealth: 100, Level: 1, Arrows: 1}, Wave: 1, Seed: seed, Status: "ready"}}
}

func Restore(state model.RunState) *Simulation {
	return &Simulation{State: model.CloneRun(state), Events: nil}
}

func (s *Simulation) Start() {
	if s.State.Status == "ready" {
		s.State.Status = "running"
		s.emit("start", s.State.ID, 0, "run started")
	}
}

func (s *Simulation) Pause() {
	if s.State.Status == "running" {
		s.State.Status = "paused"
		s.emit("pause", s.State.ID, s.State.Elapsed, "run paused")
	}
}

func (s *Simulation) Resume() {
	if s.State.Status == "paused" {
		s.State.Status = "running"
		s.emit("resume", s.State.ID, s.State.Elapsed, "run resumed")
	}
}

func (s *Simulation) Stop(reason string) {
	if s.State.Status != "finished" {
		s.State.Status = "finished"
		s.emit("stop", s.State.ID, s.State.Elapsed, reason)
	}
}

func (s *Simulation) Step(dt float64) {
	if dt <= 0 || s.State.Status != "running" {
		return
	}
	if dt > 5 {
		dt = 5
	}
	s.tick++
	s.State.Elapsed += dt
	s.spawnIfNeeded()
	s.updateEnemies(dt)
	s.autoAimAndFire(dt)
	s.updateProjectiles(dt)
	s.resolvePlayerDamage(dt)
	s.advanceLevel()
	s.State.Wave = difficulty.NextWave(s.State.Elapsed, s.State.Wave)
	if s.State.Player.Health <= 0 {
		s.Stop("player defeated")
	}
	if s.State.Elapsed >= 900 {
		s.Stop("time limit")
	}
}

func (s *Simulation) spawnIfNeeded() {
	profile := difficulty.Calculate(s.State.Elapsed, s.State.Player.Level, s.State.Wave)
	desired := difficulty.SpawnBudget(profile)
	active := model.AliveEnemyCount(s.State)
	if active >= desired {
		return
	}
	for i := active; i < desired; i++ {
		s.spawnEnemy(i, profile)
	}
}

func (s *Simulation) spawnEnemy(index int, profile difficulty.Profile) {
	kind := model.EnemyScout
	if difficulty.IsBossWave(profile.Wave) && index == 0 {
		kind = model.EnemyCaptain
	} else if index%7 == 0 {
		kind = model.EnemyBrute
	} else if index%3 == 0 {
		kind = model.EnemyWisp
	}
	health, speed, attack := difficulty.EnemyStats(string(kind), profile)
	angle := float64((index*37+int(s.State.Elapsed)*11)%360) * math.Pi / 180
	distance := 18 + float64((index*13)%9)
	enemy := model.Enemy{ID: fmt.Sprintf("%s-e-%d-%d", s.State.ID, s.tick, index), Kind: kind, Position: model.Vec2{X: math.Cos(angle) * distance, Y: math.Sin(angle) * distance}, Health: health, MaxHealth: health, Speed: speed, Attack: attack}
	s.State.Enemies = append(s.State.Enemies, enemy)
	s.emit("spawn", enemy.ID, health, string(kind))
}

func (s *Simulation) updateEnemies(dt float64) {
	for i := range s.State.Enemies {
		enemy := &s.State.Enemies[i]
		if enemy.Health <= 0 {
			continue
		}
		direction := s.State.Player.Position.Sub(enemy.Position).Normalize()
		enemy.Position = enemy.Position.Add(direction.Scale(enemy.Speed * dt))
		enemy.Age += dt
		if enemy.AttackCooldown > 0 {
			enemy.AttackCooldown -= dt
		}
	}
}

func (s *Simulation) autoAimAndFire(dt float64) {
	s.State.Player.Cooldown -= dt
	if s.State.Player.Cooldown > 0 {
		return
	}
	target, ok := s.nearestEnemy()
	if !ok {
		return
	}
	direction := target.Position.Sub(s.State.Player.Position).Normalize()
	projectile := model.Projectile{ID: fmt.Sprintf("%s-p-%d", s.State.ID, s.tick), Position: s.State.Player.Position, Velocity: direction.Scale(12), Damage: 12 + float64(s.State.Player.Level)*2, RemainingPierce: s.piercePower(), RemainingBounces: s.bouncePower(), BlastRadius: s.blastPower()}
	s.State.Projectiles = append(s.State.Projectiles, projectile)
	s.State.Player.Cooldown = s.fireCooldown()
	s.emit("fire", target.ID, projectile.Damage, "nearest target")
}

func (s *Simulation) nearestEnemy() (model.Enemy, bool) {
	var nearest model.Enemy
	best := math.MaxFloat64
	found := false
	for _, enemy := range s.State.Enemies {
		if enemy.Health <= 0 {
			continue
		}
		distance := enemy.Position.Distance(s.State.Player.Position)
		if distance < best {
			best, nearest, found = distance, enemy, true
		}
	}
	return nearest, found
}

func (s *Simulation) updateProjectiles(dt float64) {
	remaining := make([]model.Projectile, 0, len(s.State.Projectiles))
	for i := range s.State.Projectiles {
		projectile := &s.State.Projectiles[i]
		projectile.Position = projectile.Position.Add(projectile.Velocity.Scale(dt))
		projectile.Age += dt
		hit := s.hitEnemy(projectile)
		if !hit && projectile.Age < 6 {
			remaining = append(remaining, *projectile)
		}
	}
	s.State.Projectiles = remaining
}

func (s *Simulation) hitEnemy(projectile *model.Projectile) bool {
	for i := range s.State.Enemies {
		enemy := &s.State.Enemies[i]
		if enemy.Health <= 0 || projectileAlreadyHit(*projectile, enemy.ID) {
			continue
		}
		if enemy.Position.Distance(projectile.Position) > 4 {
			continue
		}
		projectile.HitIDs = append(projectile.HitIDs, enemy.ID)
		enemy.Health -= projectile.Damage
		s.emit("hit", enemy.ID, projectile.Damage, "arrow impact")
		if enemy.Health <= 0 {
			s.killEnemy(enemy, projectile)
		}
		if projectile.RemainingPierce > 0 {
			projectile.RemainingPierce--
			return false
		}
		if projectile.RemainingBounces > 0 {
			projectile.RemainingBounces--
			projectile.Velocity = projectile.Velocity.Scale(-1)
			return false
		}
		if projectile.BlastRadius > 0 {
			s.blast(enemy.Position, projectile.BlastRadius, projectile.Damage*0.5)
		}
		return true
	}
	return false
}

func projectileAlreadyHit(projectile model.Projectile, id string) bool {
	for _, hit := range projectile.HitIDs {
		if hit == id {
			return true
		}
	}
	return false
}

func (s *Simulation) blast(center model.Vec2, radius, damage float64) {
	for i := range s.State.Enemies {
		enemy := &s.State.Enemies[i]
		if enemy.Health <= 0 || enemy.Position.Distance(center) > radius*radius {
			continue
		}
		enemy.Health -= damage
		s.emit("blast", enemy.ID, damage, "area explosion")
		if enemy.Health <= 0 {
			s.killEnemy(enemy, nil)
		}
	}
}

func (s *Simulation) killEnemy(enemy *model.Enemy, projectile *model.Projectile) {
	value := 10
	if enemy.Kind == model.EnemyBrute {
		value = 18
	} else if enemy.Kind == model.EnemyCaptain {
		value = 75
	} else if enemy.Kind == model.EnemyWisp {
		value = 14
	}
	s.State.Player.Score += value
	s.State.Player.Experience += value
	if projectile != nil && model.HasUpgrade(s.State.Player, model.UpgradeOrbit) {
		s.State.Player.Score += 2
	}
	s.emit("kill", enemy.ID, float64(value), "enemy defeated")
}

func (s *Simulation) resolvePlayerDamage(dt float64) {
	for i := range s.State.Enemies {
		enemy := &s.State.Enemies[i]
		if enemy.Health <= 0 || enemy.Position.Distance(s.State.Player.Position) > 9 || enemy.AttackCooldown > 0 {
			continue
		}
		damage := enemy.Attack * dt
		if model.HasUpgrade(s.State.Player, model.UpgradeOrbit) {
			damage *= 0.8
		}
		s.State.Player.Health -= damage
		enemy.AttackCooldown = 2
		s.emit("damage", s.State.Player.ID, damage, "enemy attack")
	}
}

func (s *Simulation) advanceLevel() {
	next := difficulty.LevelForExperience(s.State.Player.Experience)
	if next <= s.State.Player.Level {
		return
	}
	for s.State.Player.Level < next {
		s.State.Player.Level++
		s.State.Player.MaxHealth += 8
		s.State.Player.Health = s.State.Player.MaxHealth
		s.emit("level", s.State.Player.ID, float64(s.State.Player.Level), "level up")
	}
}

func (s *Simulation) emit(kind, subject string, value float64, note string) {
	s.Events = append(s.Events, model.Event{Tick: s.tick, Kind: kind, Subject: subject, Value: value, Note: note})
}

func (s *Simulation) Snapshot() model.RunState {
	state := model.CloneRun(s.State)
	state.UpdatedAt = nowUTC()
	return state
}

func nowUTC() (t time.Time) { return time.Now().UTC() }
