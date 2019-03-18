//go:generate gojay -s=$GOFILE -t=User -o=generated_$GOFILE
//go:generate goimports -w  generated_$GOFILE
package user

import (
	"github.com/OGKevin/project-B-golang/interal/logging"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	// the userId
	ID uuid.UUID `gojay:"id"db:"id"`
	// the username
	Username string `gojay:"username"db:"username"`

	Password string `gojay:"-u"`
}

// CreateNewUser creates a new user.
func NewUser(username string, password []byte, users Users) (*User, error) {
	logging.Trace(logging.TraceTypeEntering)
	defer logging.Trace(logging.TraceTypeExiting)

	hash, err := hashPassword(password)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if ok, err := users.IsUsernameUnique(username); !ok {
		return nil, errors.Wrap(err, "username is not unique, or determination of uniqueness failed")
	}

	user, err := users.Create(username, hash)
	if err != nil {
		return nil, errors.Wrap(err, "could not save newly created user")
	}

	return user, nil
}

func hashPassword(password []byte) ([]byte, error) {
	logging.Trace(logging.TraceTypeEntering)
	defer logging.Trace(logging.TraceTypeExiting)

	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		logrus.WithError(err).Error("could not hash password")
		return nil, errors.Wrap(err, "could not hash password")
	}

	return hash, nil
}
