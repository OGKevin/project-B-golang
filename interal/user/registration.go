//go:generate gojay -s=$GOFILE -t=User -o=generated_$GOFILE
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
	ID uuid.UUID `json:"id"`
	// the username
	Username string `json:"username"`
}

// CreateNewUser creates a new user.
func NewUser(username, password string, users Users) (*User, error) {
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

func hashPassword(password string) ([]byte, error) {
	logging.Trace(logging.TraceTypeEntering)
	defer logging.Trace(logging.TraceTypeExiting)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithError(err).Error("could not hash password")
		return nil, errors.Wrap(err, "could not hash password")
	}

	return hash, nil
}
