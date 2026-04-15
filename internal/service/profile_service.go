package service

import (
	"context"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
	"github.com/AI-Hackathon-2026/Clients-Service/internal/repository"
)

type ProfileService interface {
	FindProfile(ctx context.Context, username string) (*model.GetUserProfilePayload, error)
	UpdateProfile(ctx context.Context, prf *model.UpdateUserProfilePayload) error
	FindGraphs(username string) ([]model.GraphPayload, error)
	GetStreak(ctx context.Context, req model.GetStreakRequest) (*model.GetStreakResponse, error)
	RegisterActivity(ctx context.Context, req model.RegisterActivityRequest) (*model.RegisterActivityResponse, error)
	AddGraph(username, graphId string) error
}

type profileService struct {
	profileRep    repository.ProfileRepository
	streakService StreakService
}

func NewProfileService(rep repository.ProfileRepository, streakService StreakService) *profileService {
	return &profileService{
		profileRep:    rep,
		streakService: streakService,
	}
}

func (s *profileService) FindProfile(ctx context.Context, username string) (*model.GetUserProfilePayload, error) {
	return s.profileRep.SelectProfileByUsername(username)
}

func (s *profileService) UpdateProfile(ctx context.Context, prf *model.UpdateUserProfilePayload) error {
	return s.profileRep.EditProfile(prf)
}

func (s *profileService) FindGraphs(username string) ([]model.GraphPayload, error) {
	return s.profileRep.SelectGraphsByUsername(username)
}

func (s *profileService) GetStreak(ctx context.Context, req model.GetStreakRequest) (*model.GetStreakResponse, error) {
	getStreakResponse, err := s.streakService.GetStreak(ctx, req)
	if err != nil {
		return nil, err
	}
	return &getStreakResponse, nil
}

func (s *profileService) RegisterActivity(ctx context.Context, req model.RegisterActivityRequest) (*model.RegisterActivityResponse, error) {
	resp, err := s.streakService.RegisterActivity(ctx, req)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *profileService) AddGraph(username, graphId string) error {
	graph := model.GraphPayload{GraphId: graphId}
	return s.profileRep.InsertGraph(username, graph)
}
