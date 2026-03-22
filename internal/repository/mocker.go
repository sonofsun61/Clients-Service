package repository

import (
	"regexp"
	"time"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
	"github.com/DATA-DOG/go-sqlmock"
)

type Mocker struct {
	mck sqlmock.Sqlmock
}

func NewMocker(mck sqlmock.Sqlmock) *Mocker {
	return &Mocker{
		mck: mck,
	}
}

func (m *Mocker) MockGetProfile() {
	clientRows := m.mck.NewRows([]string{"id", "username", "email"}).AddRow("1", "johndoe", "johndoe@mail.com")
	m.mck.ExpectQuery("SELECT \\* FROM client WHERE username=\\$1").WithArgs("johndoe").WillReturnRows(clientRows)
}

func (m *Mocker) MockEditProfile() {
	m.mck.ExpectExec(
		"UPDATE client SET username=\\$1, password=\\$2, email=\\$3 WHERE username=\\$4",
	).WithArgs(
		"johndoe1", "1234", "johndoe1@mail.com", "johndoe",
	).WillReturnResult(sqlmock.NewResult(1, 1))
}

func (m *Mocker) MockGetGraphs() {
	graphRows := m.mck.NewRows([]string{"id"}).AddRow("69b2da835143117d128a56eb")
	m.mck.ExpectQuery("SELECT graph_id FROM client c JOIN client_graph on c.id=client_id WHERE username=\\$1").WithArgs("johndoe").WillReturnRows(graphRows)
}

func (m *Mocker) MockRegister(req model.UserData) error {
	m.mck.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "user" (username, email, streak_start, last_activity) VALUES ($1, $2, $3, $4) RETURNING id`)).
		WithArgs(req.Username, req.Email, time.Now(), time.Now(), time.Now()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(req.UserID))

	m.mck.ExpectExec(regexp.QuoteMeta(`INSERT INTO auth (user_id, password) VALUES ($1, $2)`)).
		WithArgs(req.Username, req.PasswordHash).
		WillReturnResult(sqlmock.NewResult(1, 1))
	return nil
}

func (m *Mocker) MockLogin(reqEmail string) (model.UserData, error) {
	return model.UserData{}, nil
}
