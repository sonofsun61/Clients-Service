package service

import (
	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
)

type repository interface {
    SelectProfileByUsername(username string) (*model.GetUserProfilePayload, error)
    EditProfile(*model.UpdateUserProfilePayload) error
    SelectGraphsByUsername(username string) ([]model.GraphPayload, error)
}

type service struct {
    rep repository
}

func NewService(rep repository) *service {
    return &service{
        rep: rep,
    }
}

func (s *service) FindProfile(username string) (*model.GetUserProfilePayload, error) {
    return s.rep.SelectProfileByUsername(username)
}

func (s *service) UpdateProfile(prf *model.UpdateUserProfilePayload) error {
    return s.rep.EditProfile(prf)
}

func (s *service) FindGraphs(username string) ([]model.GraphPayload, error) {
    return s.rep.SelectGraphsByUsername(username) 
}
