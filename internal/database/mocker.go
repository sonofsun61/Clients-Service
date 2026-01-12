package database

import "github.com/DATA-DOG/go-sqlmock"

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
    ).WillReturnResult(sqlmock.NewResult(1,1))
}

func (m *Mocker) MockGetGraphs() {
	graphRows := m.mck.NewRows([]string{"id"}).AddRow(1).AddRow(2)
	m.mck.ExpectQuery("SELECT graph_id FROM client c JOIN client_graph on c.id=client_id WHERE username=\\$1").WithArgs("johndoe").WillReturnRows(graphRows)
}
