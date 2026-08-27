package difficulty

import "testing"

func TestDifficultyGrowsWithTimeAndLevel(t *testing.T) {
	calm := Calculate(20, 1, 1)
	hard := Calculate(300, 8, 10)
	if hard.HealthScale <= calm.HealthScale || hard.AttackScale <= calm.AttackScale || hard.SpawnRate <= calm.SpawnRate {
		t.Fatal("difficulty did not grow")
	}
}

func TestLevelThreshold(t *testing.T) {
	if LevelForExperience(80) != 2 || ExperienceForLevel(1) != 80 {
		t.Fatal("level threshold incorrect")
	}
}

func TestBossWaveAndLabel(t *testing.T) {
	if !IsBossWave(5) || IsBossWave(6) || Label(Calculate(400, 10, 5)) == "" {
		t.Fatal("wave metadata incorrect")
	}
}
