package repository

import (
	"database/sql"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
)

type ProfileRepository interface {
	SelectProfileByUsername(username string) (*model.GetUserProfilePayload, error)
	EditProfile(*model.UpdateUserProfilePayload) error
	SelectGraphsByUsername(username string) ([]model.GraphPayload, error)
}

type profileRepository struct {
	db     *sql.DB
	mocker *Mocker
}

func NewProfileRepository(db *sql.DB, mocker *Mocker) *profileRepository {
	return &profileRepository{
		db:     db,
		mocker: mocker,
	}
}

func (r *profileRepository) SelectProfileByUsername(username string) (*model.GetUserProfilePayload, error) {
	r.mocker.MockGetProfile()

	profile := new(model.GetUserProfilePayload)
	row := r.db.QueryRow("SELECT * FROM client WHERE username=$1", username)
	if err := row.Scan(&profile.UserID, &profile.Username, &profile.Email); err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *profileRepository) EditProfile(prfl *model.UpdateUserProfilePayload) error {
	r.mocker.MockEditProfile()

	_, err := r.db.Exec(
		"UPDATE client SET username=$1, password=$2, email=$3 WHERE username=$4",
		prfl.UsernameNew, prfl.Password, prfl.Email, prfl.UsernameOld,
	)
	return err
}

func (r *profileRepository) SelectGraphsByUsername(username string) ([]model.GraphPayload, error) {
	r.mocker.MockGetGraphs()

	var graphs []model.GraphPayload
	rows, err := r.db.Query("SELECT graph_id FROM client c JOIN client_graph on c.id=client_id WHERE username=$1", username)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var graph model.GraphPayload
		if err := rows.Scan(&graph.GraphId); err != nil {
			return nil, err
		}
		graphs = append(graphs, graph)
	}
	return graphs, nil
}
