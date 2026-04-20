package repository

import (
	"context"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
)

type AuthRepository interface {
	Register(ctx context.Context, userData model.UserData) error
	GetUserByEmail(ctx context.Context, reqEmail string) (model.UserData, error)
}

type authRepository struct {
	db *FakeDB
}

func NewAuthRepository(db *FakeDB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) Register(ctx context.Context, userData model.UserData) error {
	return r.db.AddUser(userData)
}

func (r *authRepository) GetUserByEmail(ctx context.Context, reqEmail string) (model.UserData, error) {
	return r.db.GetUserByEmail(reqEmail)
}
