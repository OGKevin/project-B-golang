package user

import (
	"encoding/json"
	"fmt"
	"github.com/OGKevin/project-B-golang/interal/responses"
	"github.com/asaskevich/govalidator"
	"github.com/casbin/casbin"
	"github.com/dgrijalva/jwt-go"
	"github.com/francoispqt/gojay"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type userRequest struct {
	// Username The user's username, must be unique and length(5|255)
	Username string `json:"username"valid:"length(5|255)"`
	// Password The user's password, must be length(5|255)
	Password string `json:"password"valid:"length(10|255)"`
}

type createUserRequest struct {
	user Users `gojay:"-"json:"-"`
	e *casbin.Enforcer `golay:"-"json:"-"`
}

func NewCreateUserRequest(user Users, e *casbin.Enforcer) *createUserRequest {
	return &createUserRequest{user: user, e: e}
}

// CreateUser Crates a new user
// @Summary Register a new user
// @Description Register a new user
// @ID register-new-user
// @Accept  json
// @Produce  json
// @Param body body user.userRequest true "The expected request body. Username must be length(5|255) and Password length(10|255)."
// @Success 200 {object} responses.Created "The response will include the id of the newly created user."
// @Failure 400 {object} responses.BadRequest "The error object will explain why the request failed."
// @Router /user [post]
func (b *createUserRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ur userRequest
	err := json.NewDecoder(r.Body).Decode(&ur)
	if err != nil {
		logrus.WithError(err).Error("could not decode body")
		responses.WriteBadRequests(w, responses.NewError("Request body seems to not be a valid json."))
		return
	}

	if valid, err := govalidator.ValidateStruct(ur); !valid {
		responses.WriteBadRequests(w, responses.NewValidationError(err.Error()))
		return
	}

	u, err := NewUser(ur.Username, []byte(ur.Password), b.user)
	if err != nil {
		handleError(w, r, err)
		return
	}

	 ok := b.e.AddPolicy(u.ID.String(), fmt.Sprintf("*/user/%s", u.ID), fmt.Sprintf("(%s)|(%s)|(%s)", http.MethodGet, http.MethodPut, http.MethodDelete))
	 if !ok {
	 	logrus.Error("could not add policy")
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

type getUser struct {
	// gets ID information
	users Users
}

func newGetUser(users Users) *getUser {
	return &getUser{users: users}
}

// GetUser gets user by id
// @Summary gets user by id
// @Description gets user by id
// @ID get-user
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param userId path string true "The id to get the user"
// @Param Authorization header string true "The BEARER token"
// @Success 200 {object} user.User "The user"
// @Failure 400 {object} responses.BadRequest "The error object will explain why the request failed."
// @Failure 404 {object} responses.NotFound "The error object will explain why the entity was not found."
// @Router /user/{userId} [get]
func (g getUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	userUuid, err := uuid.FromString(userId)
	if err != nil {
		responses.WriteBadRequests(w, responses.NewErrorf("%s is not a valid userID", userId))
		logrus.WithError(err).Errorf("%q is not a valid userId", userId)
		return
	}
	user, err := g.users.Get(userUuid)
	if err != nil {
		responses.WriteNotFound(w, responses.NewErrorf("user %q not found", userUuid))
		logrus.WithError(err).Errorf("user %q not found", userUuid)
		return
	}

	err = gojay.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.WithError(err).Error("Could not write response for 'get user'")
		return
	}
}

type login struct {
	users Users

	ja *jwtauth.JWTAuth
}

func newLogin(users Users, ja *jwtauth.JWTAuth) *login {
	return &login{users: users, ja: ja}
}

type jwtToken struct {
	Ack responses.Ack `json:"ack"`
	Token string `json:"token"`
}

// Login on success, you will get a JWT token to put in the auth header
// @Summary logs a user in
// @Description on success, you will get a JWT token to put in the auth header
// @ID user-login
// @Accept  json
// @Produce  json
// @Param body body user.userRequest true "The expected request body."
// @Success 200 {object} user.jwtToken "The user"
// @Failure 400 {object} responses.BadRequest "The error object will explain why the request failed."
// @Router /user/login [post]
func (l login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ur userRequest
	err := json.NewDecoder(r.Body).Decode(&ur)
	if err != nil {
		logrus.WithError(err).Error("could not decode body")
		responses.WriteBadRequests(w, responses.NewError("Request body seems to not be a valid json."))
		return
	}

	if valid, err := govalidator.ValidateStruct(ur); !valid {
		responses.WriteBadRequests(w, responses.NewValidationError(err.Error()))
		return
	}

	u, err := l.users.GetByUsername(ur.Username)
	if err != nil {
		logrus.WithError(err).Error("could not get user by username")
		http.Error(w, "username or password incorrect", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(ur.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			http.Error(w, "username or password incorrect", http.StatusUnauthorized)
			return
		}

		logrus.WithError(err).Error("could not verify user's password")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	claim := jwt.MapClaims{"user_id": u.ID.String()}
	jwtauth.SetExpiryIn(claim, time.Hour)
	jwtauth.SetIssuedNow(claim)

	t, _, err := l.ja.Encode(claim)
	if err != nil {
		logrus.WithError(err).Error("could not create new auth token")
		http.Error(w, "", http.StatusServiceUnavailable)
		return
	}

	err = json.NewEncoder(w).Encode(&jwtToken{Token:t.Raw, Ack: responses.Ack{Ack:true}})
	if err != nil {
		logrus.WithError(err).Error("could not write response")
		return
	}
}
