package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
	"github.com/AI-Hackathon-2026/Clients-Service/internal/repository"
)

type StreakService interface {
	GetStreak(ctx context.Context, req model.GetStreakRequest) (model.GetStreakResponse, error)
	RegisterActivity(ctx context.Context, req model.RegisterActivityRequest) (model.RegisterActivityResponse, error)
	ResetExpired(ctx context.Context) (int64, error)
}

type streakService struct {
	streakRepo repository.StreakRepository
	logger     *slog.Logger
	location   *time.Location
}

func NewStreakService(streakRepo repository.StreakRepository, logger *slog.Logger) StreakService {
	return &streakService{
		streakRepo: streakRepo,
		logger:     logger,
		location:   time.UTC,
	}
}

func (s *streakService) GetStreak(ctx context.Context, req model.GetStreakRequest) (model.GetStreakResponse, error) {
	state, err := s.streakRepo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return model.GetStreakResponse{}, err
	}

	today := s.nowDay()
	if state.LastActivityDate == nil {
		return model.GetStreakResponse{StreakDays: 0, IsActiveToday: false}, nil
	}

	diffDays := int(today.Sub(dayStart(*state.LastActivityDate, s.location)).Hours() / 24)
	if diffDays > 1 && state.StreakCount != 0 {
		if updateErr := s.streakRepo.UpdateByUserID(ctx, req.UserID, model.UpdateStreakRequest{
			StreakCount:      0,
			LastActivityDate: *state.LastActivityDate,
			StreakUpdatedAt:  time.Now().In(s.location),
		}); updateErr != nil {
			s.logger.Warn("failed to apply lazy streak reset", "user_id", req.UserID, "error", updateErr)
		}
		return model.GetStreakResponse{StreakDays: 0, IsActiveToday: false}, nil
	}

	return model.GetStreakResponse{
		StreakDays:    state.StreakCount,
		IsActiveToday: diffDays == 0,
	}, nil
}

func (s *streakService) RegisterActivity(ctx context.Context, req model.RegisterActivityRequest) (model.RegisterActivityResponse, error) {
	if req.UserID == "" {
		return model.RegisterActivityResponse{}, errors.New("empty user_id")
	}

	state, err := s.streakRepo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return model.RegisterActivityResponse{}, err
	}

	today := s.nowDay()
	if !req.Timestamp.IsZero() {
		today = dayStart(req.Timestamp, s.location)
	}

	nextCount := 1
	if state.LastActivityDate != nil {
		lastDay := dayStart(*state.LastActivityDate, s.location)
		diffDays := int(today.Sub(lastDay).Hours() / 24)
		switch {
		case diffDays <= 0:
			nextCount = state.StreakCount
		case diffDays == 1:
			nextCount = state.StreakCount + 1
		default:
			nextCount = 1
		}
	}

	if err := s.streakRepo.UpdateByUserID(ctx, req.UserID, model.UpdateStreakRequest{
		StreakCount:      nextCount,
		LastActivityDate: today,
		StreakUpdatedAt:  time.Now().In(s.location),
	}); err != nil {
		return model.RegisterActivityResponse{}, err
	}

	return model.RegisterActivityResponse{StreakDays: nextCount, IsActiveToday: true}, nil
}

func (s *streakService) ResetExpired(ctx context.Context) (int64, error) {
	now := time.Now().In(s.location)
	today := dayStart(now, s.location)
	return s.streakRepo.ResetExpired(ctx, today.AddDate(0, 0, -1), now)
}

func (s *streakService) nowDay() time.Time {
	return dayStart(time.Now().In(s.location), s.location)
}

func dayStart(ts time.Time, location *time.Location) time.Time {
	localTS := ts.In(location)
	y, m, d := localTS.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, location)
}
