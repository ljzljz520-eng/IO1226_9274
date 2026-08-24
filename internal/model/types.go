package model

import "time"

type Vec2 struct{ X, Y float64 }

func (v Vec2) Add(o Vec2) Vec2         { return Vec2{X: v.X + o.X, Y: v.Y + o.Y} }
func (v Vec2) Sub(o Vec2) Vec2         { return Vec2{X: v.X - o.X, Y: v.Y - o.Y} }
func (v Vec2) Scale(k float64) Vec2    { return Vec2{X: v.X * k, Y: v.Y * k} }
func (v Vec2) Distance(o Vec2) float64 { d := v.Sub(o); return (d.X*d.X + d.Y*d.Y) }
func (v Vec2) Length() float64         { return (v.X*v.X + v.Y*v.Y) }
func (v Vec2) Normalize() Vec2 {
	l := v.Length()
	if l == 0 {
		return Vec2{}
	}
	return v.Scale(1 / sqrt(l))
}

func sqrt(v float64) float64 {
	guess := v
	if guess < 1 {
		guess = 1
	}
	for i := 0; i < 12; i++ {
		guess = (guess + v/guess) / 2
	}
	return guess
}

type UpgradeKind string

const (
	UpgradePierce UpgradeKind = "pierce"
	UpgradeBounce UpgradeKind = "bounce"
	UpgradeBlast  UpgradeKind = "blast"
	UpgradeOrbit  UpgradeKind = "orbit"
)

type EnemyKind string

const (
	EnemyScout   EnemyKind = "scout"
	EnemyBrute   EnemyKind = "brute"
	EnemyWisp    EnemyKind = "wisp"
	EnemyCaptain EnemyKind = "captain"
)

type Player struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Position   Vec2          `json:"position"`
	Health     float64       `json:"health"`
	MaxHealth  float64       `json:"max_health"`
	Level      int           `json:"level"`
	Experience int           `json:"experience"`
	Score      int           `json:"score"`
	Upgrades   []UpgradeKind `json:"upgrades"`
	Arrows     int           `json:"arrows"`
	Cooldown   float64       `json:"cooldown"`
}

type Enemy struct {
	ID             string    `json:"id"`
	Kind           EnemyKind `json:"kind"`
	Position       Vec2      `json:"position"`
	Health         float64   `json:"health"`
	MaxHealth      float64   `json:"max_health"`
	Speed          float64   `json:"speed"`
	Attack         float64   `json:"attack"`
	AttackCooldown float64   `json:"attack_cooldown"`
	Age            float64   `json:"age"`
}

type Projectile struct {
	ID               string   `json:"id"`
	Position         Vec2     `json:"position"`
	Velocity         Vec2     `json:"velocity"`
	Damage           float64  `json:"damage"`
	RemainingPierce  int      `json:"remaining_pierce"`
	RemainingBounces int      `json:"remaining_bounces"`
	BlastRadius      float64  `json:"blast_radius"`
	Age              float64  `json:"age"`
	HitIDs           []string `json:"hit_ids"`
}

type RunState struct {
	ID          string       `json:"id"`
	Player      Player       `json:"player"`
	Enemies     []Enemy      `json:"enemies"`
	Projectiles []Projectile `json:"projectiles"`
	Elapsed     float64      `json:"elapsed"`
	Wave        int          `json:"wave"`
	Seed        int64        `json:"seed"`
	Status      string       `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Revision    int          `json:"revision"`
}

type Event struct {
	Tick    int     `json:"tick"`
	Kind    string  `json:"kind"`
	Subject string  `json:"subject"`
	Value   float64 `json:"value"`
	Note    string  `json:"note"`
}

type Record struct {
	ID        string    `json:"id"`
	RunID     string    `json:"run_id"`
	Sequence  int       `json:"sequence"`
	Elapsed   float64   `json:"elapsed"`
	Level     int       `json:"level"`
	Score     int       `json:"score"`
	Status    string    `json:"status"`
	Events    []Event   `json:"events"`
	CreatedAt time.Time `json:"created_at"`
}

type Batch struct {
	ID         string    `json:"id"`
	RunIDs     []string  `json:"run_ids"`
	Sequence   int       `json:"sequence"`
	State      string    `json:"state"`
	CreatedAt  time.Time `json:"created_at"`
	FinishedAt time.Time `json:"finished_at"`
}

type Audit struct {
	ID      string    `json:"id"`
	Action  string    `json:"action"`
	RunID   string    `json:"run_id"`
	BatchID string    `json:"batch_id"`
	Detail  string    `json:"detail"`
	At      time.Time `json:"at"`
}

type Profile struct {
	ID               string      `json:"id"`
	DisplayName      string      `json:"display_name"`
	BestScore        int         `json:"best_score"`
	Runs             int         `json:"runs"`
	PreferredUpgrade UpgradeKind `json:"preferred_upgrade"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

type Summary struct {
	RunID        string        `json:"run_id"`
	DisplayName  string        `json:"display_name"`
	Elapsed      float64       `json:"elapsed"`
	Level        int           `json:"level"`
	Wave         int           `json:"wave"`
	Score        int           `json:"score"`
	Kills        int           `json:"kills"`
	Status       string        `json:"status"`
	Upgrades     []UpgradeKind `json:"upgrades"`
	Difficulty   string        `json:"difficulty"`
	LatestEvents []Event       `json:"latest_events"`
}

type DispatchResult struct {
	Batch     Batch     `json:"batch"`
	Records   []Record  `json:"records"`
	Summaries []Summary `json:"summaries"`
}

type LeaderboardEntry struct {
	RunID   string
	Name    string
	Score   int
	Level   int
	Elapsed float64
}
