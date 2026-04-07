package model

import (
	"time"

	"github.com/google/uuid"
)

type GetUserProfilePayload struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UpdateUserProfilePayload struct {
	UsernameNew string `json:"username_new"`
	UsernameOld string `json:"username_old"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type GraphPayload struct {
	GraphId string `json:"graph_id"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UserData struct {
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string
	StreakStart  time.Time `json:"streak_start"`
	LastActivity time.Time `json:"last_activity"`
}

type GetStreakRequest struct {
	UserID string `json:"user_id"`
}

type GetStreakResponse struct {
	StreakDays    int  `json:"streak_days"`
	IsActiveToday bool `json:"is_active_today"`
}

type UpdateStreakRequest struct {
	StreakCount      int
	LastActivityDate time.Time
	StreakUpdatedAt  time.Time
}

type RegisterActivityRequest struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Timestamp time.Time
}

type RegisterActivityResponse struct {
	StreakDays    int  `json:"streak_days"`
	IsActiveToday bool `json:"is_active_today"`
}

type StreakState struct {
	StreakCount      int
	LastActivityDate *time.Time
	StreakUpdatedAt  time.Time
}

type StreakUpdatedEvent struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Streak    int    `json:"streak"`
	Message   string `json:"message"`
}
