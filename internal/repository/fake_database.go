package repository

import (
	"errors"
	"time"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
)

type FakeDB struct {
	users      map[string]model.UserData
	userGraphs map[string][]model.GraphPayload
}

func NewFakeDB() *FakeDB {
	return &FakeDB{
		users:      make(map[string]model.UserData),
		userGraphs: make(map[string][]model.GraphPayload),
	}
}

func (db *FakeDB) AddUser(userData model.UserData) error {
	for _, user := range db.users {
		if user.Email == userData.Email {
			return errors.New("user with this email already exists")
		}
		if user.Username == userData.Username {
			return errors.New("user with this username already exists")
		}
	}

	db.users[userData.UserID.String()] = userData
	db.userGraphs[userData.UserID.String()] = []model.GraphPayload{}

	return nil
}

func (db *FakeDB) GetUserByEmail(email string) (model.UserData, error) {
	for _, user := range db.users {
		if user.Email == email {
			return user, nil
		}
	}
	return model.UserData{}, errors.New("user not found")
}

func (db *FakeDB) GetUserByUsername(username string) (model.UserData, error) {
	for _, user := range db.users {
		if user.Username == username {
			return user, nil
		}
	}
	return model.UserData{}, errors.New("user not found")
}

func (db *FakeDB) GetProfileByUsername(username string) (*model.GetUserProfilePayload, error) {
	user, err := db.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return &model.GetUserProfilePayload{
		UserID:   user.UserID.String(),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (db *FakeDB) UpdateProfile(prfl *model.UpdateUserProfilePayload) error {
	user, err := db.GetUserByUsername(prfl.UsernameOld)
	if err != nil {
		return err
	}

	if prfl.UsernameNew != "" && prfl.UsernameNew != prfl.UsernameOld {
		for _, u := range db.users {
			if u.Username == prfl.UsernameNew && u.UserID != user.UserID {
				return errors.New("username already taken")
			}
		}
		user.Username = prfl.UsernameNew
	}

	if prfl.Email != "" && prfl.Email != user.Email {
		for _, u := range db.users {
			if u.Email == prfl.Email && u.UserID != user.UserID {
				return errors.New("email already taken")
			}
		}
		user.Email = prfl.Email
	}

	db.users[user.UserID.String()] = user

	return nil
}

func (db *FakeDB) GetGraphsByUsername(username string) ([]model.GraphPayload, error) {
	user, err := db.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	graphs, exists := db.userGraphs[user.UserID.String()]
	if !exists {
		return []model.GraphPayload{}, nil
	}

	return graphs, nil
}

func (db *FakeDB) AddGraphToUser(userID string, graph model.GraphPayload) error {
	if _, exists := db.users[userID]; !exists {
		return errors.New("user not found")
	}

	db.userGraphs[userID] = append(db.userGraphs[userID], graph)
	return nil
}

func (db *FakeDB) GetStreakByUserID(userID string) (model.StreakState, error) {
	user, exists := db.users[userID]
	if !exists {
		return model.StreakState{}, errors.New("user not found")
	}

	return model.StreakState{
		StreakCount:      user.StreakStart.Day() - user.LastActivity.Day(),
		LastActivityDate: &user.LastActivity,
	}, nil
}

func (db *FakeDB) UpdateStreakByUserID(userID string, req model.UpdateStreakRequest) error {
	user, exists := db.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	if !req.LastActivityDate.IsZero() {
		user.LastActivity = req.LastActivityDate
	}
	user.StreakStart = req.StreakUpdatedAt

	db.users[userID] = user
	return nil
}

func (db *FakeDB) ResetExpiredStreaks(cutoffDate time.Time, now time.Time) (int64, error) {
	var resetCount int64 = 0

	for userID, user := range db.users {
		if user.LastActivity.Before(cutoffDate) {
			user.StreakStart = now
			user.LastActivity = now
			db.users[userID] = user
			resetCount++
		}
	}

	return resetCount, nil
}

func (db *FakeDB) GetUserByID(userID string) (model.UserData, error) {
	user, exists := db.users[userID]
	if !exists {
		return model.UserData{}, errors.New("user not found")
	}

	return user, nil
}

func (db *FakeDB) GetAllUsers() []model.UserData {
	users := make([]model.UserData, 0, len(db.users))
	for _, user := range db.users {
		users = append(users, user)
	}
	return users
}

func (db *FakeDB) DeleteUser(userID string) error {
	_, exists := db.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	delete(db.users, userID)
	delete(db.userGraphs, userID)

	return nil
}

func (db *FakeDB) UpdateUserPassword(userID string, newPasswordHash string) error {
	user, exists := db.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	user.PasswordHash = newPasswordHash
	db.users[userID] = user

	return nil
}

func (db *FakeDB) GetStreakByUserIDWithReset(userID string, currentDate time.Time, resetThresholdHours int) (model.GetStreakResponse, error) {
	user, err := db.GetUserByID(userID)
	if err != nil {
		return model.GetStreakResponse{}, err
	}

	isActiveToday := false
	lastActivityDate := user.LastActivity
	daysDiff := int(currentDate.Sub(lastActivityDate).Hours() / 24)
	isActiveToday = daysDiff == 0

	streakDays := int(currentDate.Sub(user.StreakStart).Hours() / 24)
	if streakDays < 0 {
		streakDays = 0
	}

	return model.GetStreakResponse{
		StreakDays:    streakDays,
		IsActiveToday: isActiveToday,
	}, nil
}

func (db *FakeDB) RegisterUserActivity(userID string, sessionID string, timestamp time.Time) (model.RegisterActivityResponse, error) {
	user, err := db.GetUserByID(userID)
	if err != nil {
		return model.RegisterActivityResponse{}, err
	}

	lastActivity := user.LastActivity
	streakDays := int(timestamp.Sub(user.StreakStart).Hours() / 24)
	if streakDays < 0 {
		streakDays = 0
	}
	isActiveToday := false

	lastDateTrunc := time.Date(lastActivity.Year(), lastActivity.Month(), lastActivity.Day(), 0, 0, 0, 0, lastActivity.Location())
	todayTrunc := time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, timestamp.Location())

	daysDiff := int(todayTrunc.Sub(lastDateTrunc).Hours() / 24)

	if daysDiff == 0 {
		isActiveToday = true
	} else if daysDiff == 1 {
		streakDays++
		isActiveToday = true
		user.StreakStart = timestamp
	} else if daysDiff > 1 {
		streakDays = 1
		isActiveToday = true
		user.StreakStart = timestamp
	}

	user.LastActivity = timestamp
	db.users[userID] = user

	return model.RegisterActivityResponse{
		StreakDays:    streakDays,
		IsActiveToday: isActiveToday,
	}, nil
}
