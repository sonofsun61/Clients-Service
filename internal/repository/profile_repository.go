package repository

import (
	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
)

type ProfileRepository interface {
	SelectProfileByUsername(username string) (*model.GetUserProfilePayload, error)
	EditProfile(*model.UpdateUserProfilePayload) error
	SelectGraphsByUsername(username string) ([]model.GraphPayload, error)
	InsertGraph(username string, graphId model.GraphPayload) error
}

type profileRepository struct {
	db *FakeDB
}

func NewProfileRepository(db *FakeDB) *profileRepository {
	return &profileRepository{
		db: db,
	}
}

func (r *profileRepository) SelectProfileByUsername(username string) (*model.GetUserProfilePayload, error) {
	return r.db.GetProfileByUsername(username)
}

func (r *profileRepository) EditProfile(prfl *model.UpdateUserProfilePayload) error {
	return r.db.UpdateProfile(prfl)
}

func (r *profileRepository) SelectGraphsByUsername(username string) ([]model.GraphPayload, error) {
	return r.db.GetGraphsByUsername(username)
}

func (r *profileRepository) InsertGraph(username string, graphId model.GraphPayload) error {
	return r.db.AddGraphToUser(username, graphId)
}
