package model

import (
	"time"

	"github.com/google/uuid"
)

type GetUserProfilePayload struct {
	Id       string `json:"id"`
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
