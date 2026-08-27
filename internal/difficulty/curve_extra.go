package difficulty

func ThreatAt(elapsed float64, level int) float64 {
	p := Calculate(elapsed, level, 1)
	return p.HealthScale*p.AttackScale + p.SpeedScale
}

func IsBossWave(wave int) bool { return wave > 0 && wave%5 == 0 }

func RecommendedLevel(elapsed float64) int {
	level := int(elapsed/45) + 1
	if level > 20 {
		return 20
	}
	return level
}

func RecoveryWindow(profile Profile) float64 {
	window := 8.0 / profile.SpawnRate
	if window < 1.5 {
		return 1.5
	}
	return window
}

func ScoreMultiplier(profile Profile) float64 {
	value := 1 + profile.HealthScale*0.25 + profile.SpeedScale*0.1
	if IsBossWave(profile.Wave) {
		value += 0.8
	}
	return value
}
