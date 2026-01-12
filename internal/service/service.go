package service

import (
	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
)

type repository interface {
    selectProfileByUsername(username string) (*model.GetUserProfilePayload, error)
    editProfile(prf *model.UpdateUserProfilePayload) error
    selectGraphsByUsername(username string) ([]string, error)
}

type service struct {
    rep repository
}

func NewService(rep repository) *service {
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

func (s *service) findGraphs(username string) ([]string, error) {
    return s.rep.selectGraphsByUsername(username) 
}
