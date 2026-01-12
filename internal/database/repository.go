package database

import (
	"database/sql"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
)

type repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
    return &repository{
        db: db,
    }
}

func (r *repository) selectProfileByUsername(username string) (*model.GetUserProfilePayload, error) {
    profile := new(model.GetUserProfilePayload)
    row := r.db.QueryRow("SELECT * FROM client WHERE username=$1", username)
    if err := row.Scan(&profile.Id, &profile.Username, &profile.Email); err != nil {
        return nil, err
    }
    return profile, nil
}

func (r *repository) editProfile(prfl *model.UpdateUserProfilePayload) error {
    _, err := r.db.Exec(
        "UPDATE client SET username=$1, password=$2, email=$3 WHERE username=$4",
        prfl.UsernameNew, prfl.Password, prfl.Email, prfl.UsernameOld,
    )
    return err
}

func (r *repository) selectGraphsByUsername(username string) ([]string, error) {
    var graphsId []string
    rows, err := r.db.Query("SELECT graph_id FROM client c JOIN client_graph on c.id = client_id WHERE username=$1", username)
    if err != nil {
        return nil, err
    }
    for rows.Next() {
        var id string
        if err := rows.Scan(&id); err != nil {
            return nil, err
        }
        graphsId = append(graphsId, id)
    }
    return graphsId, nil
}
