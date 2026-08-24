package engine

import "miniarrow/internal/model"

type Replay struct {
	Initial model.RunState
	Frames  []model.RunState
}

func NewReplay(initial model.RunState) *Replay {
	return &Replay{Initial: model.CloneRun(initial), Frames: make([]model.RunState, 0)}
}

func (r *Replay) Capture(state model.RunState) { r.Frames = append(r.Frames, model.CloneRun(state)) }

func (r *Replay) Last() model.RunState {
	if len(r.Frames) == 0 {
		return model.CloneRun(r.Initial)
	}
	return model.CloneRun(r.Frames[len(r.Frames)-1])
}

func (r *Replay) Frame(index int) (model.RunState, bool) {
	if index < 0 || index >= len(r.Frames) {
		return model.RunState{}, false
	}
	return model.CloneRun(r.Frames[index]), true
}

func (r *Replay) Duration() float64 {
	if len(r.Frames) == 0 {
		return 0
	}
	return r.Frames[len(r.Frames)-1].Elapsed - r.Initial.Elapsed
}

func (r *Replay) AliveAt(index int) int {
	frame, ok := r.Frame(index)
	if !ok {
		return 0
	}
	return model.AliveEnemyCount(frame)
}

func (r *Replay) Summary() map[string]any {
	last := r.Last()
	return map[string]any{"frames": len(r.Frames), "duration": r.Duration(), "status": last.Status, "score": last.Player.Score, "level": last.Player.Level}
}
