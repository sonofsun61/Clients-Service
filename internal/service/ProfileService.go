package service

import (
	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
)

type ProfileRepository interface {
	SelectProfileByUsername(username string) (*model.GetUserProfilePayload, error)
	EditProfile(*model.UpdateUserProfilePayload) error
	SelectGraphsByUsername(username string) ([]model.GraphPayload, error)
}

type ProfileService struct {
	rep ProfileRepository
}

func NewProfileService(rep ProfileRepository) *ProfileService {
	return &ProfileService{
		rep: rep,
	}
}

func (s *ProfileService) FindProfile(username string) (*model.GetUserProfilePayload, error) {
	return s.rep.SelectProfileByUsername(username)
}

func (s *ProfileService) UpdateProfile(prf *model.UpdateUserProfilePayload) error {
	return s.rep.EditProfile(prf)
}

func (s *ProfileService) FindGraphs(username string) ([]model.GraphPayload, error) {
	return s.rep.SelectGraphsByUsername(username)
}
