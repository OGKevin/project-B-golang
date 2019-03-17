//go:generate gojay -s=$GOFILE -t=createUserRequest -o=generated_$GOFILE
package user

import (
	"github.com/OGKevin/project-B-golang/interal/responses"
	"github.com/asaskevich/govalidator"
	"github.com/francoispqt/gojay"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

type createUserRequest struct {
	Username string `gojay:"username"valid:"length(5|255)"`
	Password string `gojay:"password"valid:"length(10|255)"`

	user Users `gojay:"-"json:"-"`
}

func NewCreateUserRequest(user Users) *createUserRequest {
	return &createUserRequest{user: user}
}

// CreateUser Crates a new user
// @Summary Register a new user
// @Description Register a new user
// @ID register-new-user
// @Accept  json
// @Produce  json
// @Param body body user.createUserRequest true "The expected request body. Username must be length(5|255) and Password length(10|255)."
// @Success 200 {object} responses.Created "The response will include the id of the newly created user."
// @Failure 400 {object} responses.BadRequest "The error object will explain why the request failed."
// @Router /user [post]
func (b *createUserRequest) ServeHttp(w http.ResponseWriter, r *http.Request) {
	err := gojay.NewDecoder(r.Body).DecodeObject(b)
	if err != nil {
		logrus.WithError(err).Error("could not decode body")
		responses.WriteBadRequests(w, responses.NewError("Request body seems to not be a valid json."))
		return
	}

	if valid, err := govalidator.ValidateStruct(b); !valid {
		responses.WriteBadRequests(w, responses.NewValidationError(err.Error()))
		return
	}

	u, err := NewUser(b.Username, b.Password, b.user)
	if err != nil {
		handleError(w, r, err)
		return
	}

	responses.WriteCreated(w, u.ID)
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	type causer interface {
		Cause() error
	}

	logrus.WithError(err).Error("error occurred while handling request")

	switch err.(type) {
	case *usernameNotUnique:
		responses.WriteBadRequests(w, responses.NewError(err.Error()))
	case causer:
		handleError(w, r, errors.Cause(err))
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
