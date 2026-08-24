package difficulty

import "math"

type Tuning struct {
	BaseSpawn   float64
	TimeWeight  float64
	LevelWeight float64
	WaveWeight  float64
	BossBonus   float64
}

func DefaultTuning() Tuning {
	return Tuning{BaseSpawn: 1, TimeWeight: 0.045, LevelWeight: 0.06, WaveWeight: 0.025, BossBonus: 0.8}
}

func NormalizeTuning(t Tuning) Tuning {
	if t.BaseSpawn <= 0 {
		t.BaseSpawn = 1
	}
	if t.TimeWeight < 0 {
		t.TimeWeight = 0
	}
	if t.LevelWeight < 0 {
		t.LevelWeight = 0
	}
	if t.WaveWeight < 0 {
		t.WaveWeight = 0
	}
	if t.BossBonus < 0 {
		t.BossBonus = 0
	}
	return t
}

func Evaluate(t Tuning, elapsed float64, level, wave int) Profile {
	t = NormalizeTuning(t)
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
	p := Profile{Elapsed: elapsed, Level: level, Wave: wave}
	p.SpawnRate = t.BaseSpawn + timeFactor*t.TimeWeight + float64(level-1)*t.LevelWeight + float64(wave-1)*t.WaveWeight
	p.HealthScale = 1 + timeFactor*0.12 + float64(level-1)*0.18 + float64(wave-1)*0.07
	p.SpeedScale = 1 + timeFactor*0.035 + float64(level-1)*0.04 + float64(wave-1)*0.02
	p.AttackScale = 1 + timeFactor*0.08 + float64(level-1)*0.1 + float64(wave-1)*0.05
	p.Label = band(p.HealthScale)
	return p
}

func band(health float64) string {
	if health >= 3.5 {
		return "nightmare"
	}
	if health >= 2.2 {
		return "veteran"
	}
	if health >= 1.4 {
		return "hardened"
	}
	return "calm"
}

func SpawnSchedule(elapsed float64, wave int) []float64 {
	profile := Calculate(elapsed, 1, wave)
	count := SpawnBudget(profile)
	out := make([]float64, count)
	interval := 10 / profile.SpawnRate
	for i := range out {
		out[i] = math.Round(float64(i)*interval*100) / 100
	}
	return out
}

func WaveDuration(wave int) float64 {
	if wave < 1 {
		wave = 1
	}
	return 25 + float64(wave%4)*5
}

func IsOverwhelmed(profile Profile, active int) bool { return active > SpawnBudget(profile)*2 }

func RecoveryScore(profile Profile, health float64) float64 {
	if health <= 0 {
		return 0
	}
	score := health / (profile.HealthScale * profile.AttackScale)
	if score > 100 {
		score = 100
	}
	return score
}
