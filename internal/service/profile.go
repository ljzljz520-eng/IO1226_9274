package service

import (
	"errors"
	"miniarrow/internal/model"
	"time"
)

func (a *App) UpdateProfile(id, displayName string, preferred model.UpgradeKind) (model.Profile, error) {
	if id == "" || displayName == "" {
		return model.Profile{}, errors.New("profile id and display name are required")
	}
	profile, err := a.store.GetProfile(id)
	if err != nil {
		profile = model.Profile{ID: id}
	}
	profile.DisplayName = displayName
	profile.PreferredUpgrade = preferred
	profile.UpdatedAt = time.Now().UTC()
	return profile, a.store.PutProfile(profile)
}

func (a *App) Profile(id string) (model.Profile, error) { return a.store.GetProfile(id) }

func (a *App) SaveBest(profileID string, summary model.Summary) error {
	profile, err := a.store.GetProfile(profileID)
	if err != nil {
		return err
	}
	if summary.Score > profile.BestScore {
		profile.BestScore = summary.Score
	}
	profile.Runs++
	profile.UpdatedAt = time.Now().UTC()
	return a.store.PutProfile(profile)
}
