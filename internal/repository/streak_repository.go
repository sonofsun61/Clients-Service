package repository

import (
	"context"
	"time"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
)

type StreakRepository interface {
	GetByUserID(ctx context.Context, userID string) (model.StreakState, error)
	UpdateByUserID(ctx context.Context, userID string, req model.UpdateStreakRequest) error
	ResetExpired(ctx context.Context, cutoffDate time.Time, now time.Time) (int64, error)
}

type streakRepository struct {
	db *FakeDB
}

func NewStreakRepository(db *FakeDB) StreakRepository {
	return &streakRepository{db: db}
}

func (r *streakRepository) GetByUserID(ctx context.Context, userID string) (model.StreakState, error) {
	return r.db.GetStreakByUserID(userID)
}

func (r *streakRepository) UpdateByUserID(ctx context.Context, userID string, req model.UpdateStreakRequest) error {
	return r.db.UpdateStreakByUserID(userID, req)
}

func (r *streakRepository) ResetExpired(ctx context.Context, cutoffDate time.Time, now time.Time) (int64, error) {
	return r.db.ResetExpiredStreaks(cutoffDate, now)
}
