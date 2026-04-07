package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
)

type StreakRepository interface {
	GetByUserID(ctx context.Context, userID string) (model.StreakState, error)
	UpdateByUserID(ctx context.Context, userID string, req model.UpdateStreakRequest) error
	ResetExpired(ctx context.Context, cutoffDate time.Time, now time.Time) (int64, error)
}

type streakRepository struct {
	db *sql.DB
}

func NewStreakRepository(db *sql.DB) StreakRepository {
	return &streakRepository{db: db}
}

func (r *streakRepository) GetByUserID(ctx context.Context, userID string) (model.StreakState, error) {
	var state model.StreakState
	var lastActivityDate sql.NullTime

	row := r.db.QueryRowContext(
		ctx,
		`SELECT COALESCE(streak_count, 0), last_activity_date, COALESCE(streak_updated_at, NOW()) FROM client WHERE id = $1`,
		userID,
	)

	if err := row.Scan(&state.StreakCount, &lastActivityDate, &state.StreakUpdatedAt); err != nil {
		return model.StreakState{}, err
	}
	if lastActivityDate.Valid {
		last := lastActivityDate.Time
		state.LastActivityDate = &last
	}
	return state, nil
}

func (r *streakRepository) UpdateByUserID(ctx context.Context, userID string, req model.UpdateStreakRequest) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE client SET streak_count = $1, last_activity_date = $2, streak_updated_at = $3 WHERE id = $4`,
		req.StreakCount,
		req.LastActivityDate,
		req.StreakUpdatedAt,
		userID,
	)
	return err
}

func (r *streakRepository) ResetExpired(ctx context.Context, cutoffDate time.Time, now time.Time) (int64, error) {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE client
		 SET streak_count = 0, streak_updated_at = $1
		 WHERE streak_count > 0 AND last_activity_date IS NOT NULL AND last_activity_date < $2`,
		now,
		cutoffDate,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
