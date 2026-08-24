package model

func ClonePlayer(in Player) Player {
	out := in
	out.Upgrades = append([]UpgradeKind(nil), in.Upgrades...)
	return out
}

func CloneEnemy(in Enemy) Enemy { return in }

func CloneProjectile(in Projectile) Projectile {
	out := in
	out.HitIDs = append([]string(nil), in.HitIDs...)
	return out
}

func CloneRun(in RunState) RunState {
	out := in
	out.Player = ClonePlayer(in.Player)
	out.Enemies = make([]Enemy, len(in.Enemies))
	for i, enemy := range in.Enemies {
		out.Enemies[i] = CloneEnemy(enemy)
	}
	out.Projectiles = make([]Projectile, len(in.Projectiles))
	for i, projectile := range in.Projectiles {
		out.Projectiles[i] = CloneProjectile(projectile)
	}
	return out
}

func CloneEvents(in []Event) []Event {
	out := make([]Event, len(in))
	copy(out, in)
	return out
}

func AddUpgrade(player *Player, kind UpgradeKind) {
	for _, existing := range player.Upgrades {
		if existing == kind {
			return
		}
	}
	player.Upgrades = append(player.Upgrades, kind)
}

func HasUpgrade(player Player, kind UpgradeKind) bool {
	for _, existing := range player.Upgrades {
		if existing == kind {
			return true
		}
	}
	return false
}

func UpgradeCount(player Player) int { return len(player.Upgrades) }

func EnemyCount(run RunState) int { return len(run.Enemies) }

func AliveEnemyCount(run RunState) int {
	count := 0
	for _, enemy := range run.Enemies {
		if enemy.Health > 0 {
			count++
		}
	}
	return count
}
