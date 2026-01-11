package service

import (
	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
)

type Repository interface {
    selectProfileByUsername(username string) (*model.GetUserProfilePayload, error)
    editProfile(*model.UpdateUserProfilePayload) error
    selectGraphsByUsername(username string) ([]model.GraphPayload, error)
}

type service struct {
    rep Repository
}

func NewService(rep Repository) *service {
    return &service{
        rep: rep,
    }
}

func (s *service) findProfile(username string) (*model.GetUserProfilePayload, error) {
    return s.rep.selectProfileByUsername(username)
}

func (s *service) updateProfile(prf *model.UpdateUserProfilePayload) error {
    return s.rep.editProfile(prf)
}

func (s *service) findGraphs(username string) ([]model.GraphPayload, error) {
    return s.rep.selectGraphsByUsername(username) 
}
