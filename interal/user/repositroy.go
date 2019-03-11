package user

import (
	"github.com/OGKevin/project-B-golang/interal/database"
	"github.com/OGKevin/project-B-golang/interal/logging"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type Users interface {
	Create(username string, password []byte) (*User, error)
	IsUsernameUnique(username string) (bool, error)
}

type UsersDatabase struct {
	db *sqlx.DB
}

func (u *UsersDatabase) IsUsernameUnique(username string) (bool, error) {
	rows, err := u.db.Query(`select count(*) from users where username = ?`, username)
	if err != nil {
		return false, errors.Wrap(err, "could not determine if user name is unique")
	}

	var c int
	if !rows.Next() {
		return false, errors.New("there seems to be no results in rows")
	}

	err = rows.Scan(&c)
	if err != nil {
		return false, errors.Wrap(err, "could not scan rows to determine if username is unique")
	}

	if c != 0 {
		return false, &usernameNotUnique{username: username}
	}

	return true, nil
}

// NewUsersDatabase creates a new repo that talks to the db.
func NewUsersDatabase(db *sqlx.DB) *UsersDatabase {
	var d *sqlx.DB
	if db == nil {
		d = database.GetDB()
	} else {
		d = db
	}

	return &UsersDatabase{db: d}
}

func (u *UsersDatabase) Create(username string, password []byte) (*User, error) {
	logging.Trace(logging.TraceTypeEntering)
	defer logging.Trace(logging.TraceTypeExiting)

	id := uuid.NewV4()
	_, err := u.db.Exec(
		`insert into users (id, username, password) value (?, ?, ?)`,
		id, username, password,
		)
	if err != nil {
		return nil, errors.Wrap(err, "could not create new user record")
	}

	return &User{ID: id, Username: username}, nil
}
