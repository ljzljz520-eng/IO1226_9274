package difficulty

import "fmt"

type Profile struct {
	Elapsed     float64
	Level       int
	Wave        int
	SpawnRate   float64
	HealthScale float64
	SpeedScale  float64
	AttackScale float64
	Label       string
}

func Calculate(elapsed float64, level, wave int) Profile {
	if elapsed < 0 {
		elapsed = 0
	}
	if level < 1 {
		level = 1
	}
	if wave < 1 {
		wave = 1
	}
	timeFactor := elapsed / 60
	levelFactor := float64(level - 1)
	waveFactor := float64(wave - 1)
	profile := Profile{Elapsed: elapsed, Level: level, Wave: wave}
	profile.SpawnRate = 1.0 + timeFactor*0.045 + levelFactor*0.06 + waveFactor*0.025
	profile.HealthScale = 1.0 + timeFactor*0.12 + levelFactor*0.18 + waveFactor*0.07
	profile.SpeedScale = 1.0 + timeFactor*0.035 + levelFactor*0.04 + waveFactor*0.02
	profile.AttackScale = 1.0 + timeFactor*0.08 + levelFactor*0.1 + waveFactor*0.05
	switch {
	case profile.HealthScale >= 3.5:
		profile.Label = "nightmare"
	case profile.HealthScale >= 2.2:
		profile.Label = "veteran"
	case profile.HealthScale >= 1.4:
		profile.Label = "hardened"
	default:
		profile.Label = "calm"
	}
	return profile
}

func SpawnBudget(profile Profile) int {
	budget := int(profile.SpawnRate*3) + profile.Wave
	if budget < 3 {
		return 3
	}
	if budget > 40 {
		return 40
	}
	return budget
}

func EnemyStats(kind string, profile Profile) (health, speed, attack float64) {
	baseHealth, baseSpeed, baseAttack := 20.0, 1.0, 3.0
	switch kind {
	case "brute":
		baseHealth, baseSpeed, baseAttack = 55, 0.55, 8
	case "wisp":
		baseHealth, baseSpeed, baseAttack = 14, 1.7, 2
	case "captain":
		baseHealth, baseSpeed, baseAttack = 110, 0.8, 14
	}
	return baseHealth * profile.HealthScale, baseSpeed * profile.SpeedScale, baseAttack * profile.AttackScale
}

func Label(profile Profile) string { return fmt.Sprintf("%s-%d", profile.Label, profile.Wave) }

func NextWave(elapsed float64, current int) int {
	if current < 1 {
		current = 1
	}
	threshold := float64(current) * 25
	if elapsed >= threshold {
		return current + 1
	}
	return current
}

func ExperienceForLevel(level int) int {
	if level < 1 {
		level = 1
	}
	return 80 + (level-1)*45 + (level-1)*(level-1)*10
}

func LevelForExperience(experience int) int {
	level := 1
	for experience >= ExperienceForLevel(level) {
		level++
	}
	return level
}
