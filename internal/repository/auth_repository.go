package repository

import (
	"context"
	"database/sql"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
)

type AuthRepository interface {
	Register(ctx context.Context, userData model.UserData) error
	GetUserByEmail(ctx context.Context, reqEmail string) (model.UserData, error)
}

type authRepository struct {
	db     *sql.DB
	mocker *Mocker
}

func NewAuthRepository(db *sql.DB, mocker *Mocker) AuthRepository {
	return &authRepository{
		db:     db,
		mocker: mocker,
	}
}

func (r *authRepository) Register(ctx context.Context, userData model.UserData) error {
	return r.mocker.MockRegister(userData)
	// query := `
	// 	INSERT INTO users (user_id, username, password_hash, streak_start, last_activity)
	// 	VALUES ($1, $2, $3, $4, $5)
	// `
	// _, err := r.db.ExecContext(ctx, query, userData.UserID, userData.Username, userData.PasswordHash, time.Now(), time.Now())
	// if err != nil {
	// 	return fmt.Errorf("failed to insert user: %w", err)
	// }
	// return nil
}

func (r *authRepository) GetUserByEmail(ctx context.Context, reqEmail string) (model.UserData, error) {
	return r.mocker.MockLogin(reqEmail)
	// query := `
	// 	SELECT user_id, username, email, password_hash, streak_start, last_activity
	// 	FROM users
	// 	WHERE email = $1
	// `
	// var userData model.UserData
	// err := r.db.QueryRowContext(ctx, query, reqEmail).Scan(
	// 	&userData.UserID,
	// 	&userData.Username,
	// 	&userData.Email,
	// 	&userData.PasswordHash,
	// 	&userData.StreakStart,
	// 	&userData.LastActivity,
	// )
	// if err != nil {
	// 	return model.UserData{}, fmt.Errorf("query failed: %w", err)
	// }
	// return userData, nil
}
