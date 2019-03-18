package user

import (
	"bytes"
	"fmt"
	"github.com/OGKevin/project-B-golang/interal/acl"
	"github.com/OGKevin/xorm-adapter"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/francoispqt/gojay"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func newResponseWriter() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

func TestCreateUser(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name  string
		args  args
		code  int
		users Users
	}{
		{
			name: "created",
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodPost, "/", createUserBody(t)),
			},
			code:  201,
			users: &usersForTest{},
		},
		{
			name: "username not unique",
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodPost, "/", createUserBody(t)),
			},
			code:  400,
			users: &usersForTest{usernameNotUnique: true},
		},
		{
			name: "user creation failed",
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodPost, "/", createUserBody(t)),
			},
			code:  500,
			users: &usersForTest{creationFailed: true},
		},
		{
			name: "no body",
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("")),
			},
			code:  400,
			users: &usersForTest{},
		},
		{
			name: "bad body",
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{}")),
			},
			code:  400,
			users: &usersForTest{},
		},
	}

	ja := jwtauth.New("HS256", []byte("secret"), nil)
	e := acl.NewEnforcer(xormadapter.NewAdapter("sqlite3", "file::memory:?mode=memory&cache=shared"))
	e.EnableLog(true)
	e.EnableAutoSave(true)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := BuildRouter(tt.users, ja, e)
			m.ServeHTTP(tt.args.w, tt.args.r)

			w := tt.args.w.(*httptest.ResponseRecorder)
			if !assert.Equal(t, tt.code, w.Code) {
				return
			}

			if w.Code == http.StatusCreated {
				var u User
				if !assert.NoError(t, gojay.NewDecoder(w.Body).DecodeObject(&u)) {
					return
				}

				w := newResponseWriter()
				r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", u.ID), nil)
				token, _, err := ja.Encode(jwt.MapClaims{"user_id": u.ID.String()})
				if !assert.NoError(t, err) {
					return
				}

				r.Header.Set("Authorization", fmt.Sprintf("BEARER %s", token.Raw))

				m.ServeHTTP(w, r)

				assert.Equal(t, http.StatusOK, w.Code)
			}
		})
	}
}

func createUserBody(t *testing.T) io.Reader {
	buf := bytes.NewBufferString("")
	err := gojay.NewEncoder(buf).EncodeObject(&createUserRequest{Username: uuid.NewV4().String(), Password: uuid.NewV4().String()})
	assert.NoError(t, err)

	return buf
}

type usersForTest struct {
	usernameNotUnique bool
	creationFailed    bool

	userNotFound bool
}

func (u *usersForTest) Get(id uuid.UUID) (*User, error) {
	if u.userNotFound {
		return nil, errors.New("user not found")
	}

	return u.Create("Sjaak", nil)
}

func (u *usersForTest) IsUsernameUnique(username string) (bool, error) {
	if u.usernameNotUnique {
		return false, &usernameNotUnique{username: username}
	}
	return true, nil
}

func (u *usersForTest) Create(username string, password []byte) (*User, error) {
	if u.creationFailed {
		return nil, errors.New("creation failed")
	}

	return &User{ID: uuid.NewV4(), Username: username}, nil
}

func Test_getUser_ServeHTTP(t *testing.T) {
	type fields struct {
		users Users
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		code   int
		userID uuid.UUID
		disableACL bool
	}{
		{
			name: "get user",
			fields: fields{
				users: &usersForTest{},
			},
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", uuid.FromStringOrNil("F852243B-6BE8-4AB3-A557-4584CF19ABB0")), nil),
			},
			code: http.StatusOK,
			userID: uuid.FromStringOrNil("F852243B-6BE8-4AB3-A557-4584CF19ABB0"),
		},
		{
			name: "get user with invalid id",
			fields: fields{
				users: &usersForTest{},
			},
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodGet, "/jhfsdkfhgasdjklfhjsd", nil),
			},
			code: http.StatusForbidden,
			userID: uuid.NewV4(),
		},
		{
			name: "get user with invalid id w o acl",
			fields: fields{
				users: &usersForTest{},
			},
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodGet, "/jhfsdkfhgasdjklfhjsd", nil),
			},
			code: http.StatusBadRequest,
			userID: uuid.NewV4(),
			disableACL: true,
		},
		{
			name: "get user id not found w o acl",
			fields: fields{
				users: &usersForTest{userNotFound: true},
			},
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", uuid.NewV4()), nil),
			},
			code: http.StatusNotFound,
			userID: uuid.NewV4(),
			disableACL: true,

		},
		{
			name: "get user with valid id not found",
			fields: fields{
				users: &usersForTest{userNotFound: true},
			},
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", uuid.FromStringOrNil("AC0E1E53-5D16-4F3C-8780-54C3BBE02CDB")), nil),
			},
			code: http.StatusNotFound,
			userID: uuid.FromStringOrNil("AC0E1E53-5D16-4F3C-8780-54C3BBE02CDB"),
		},
	}

	ja := jwtauth.New("HS256", []byte("secret"), nil)
	e := acl.NewEnforcer(xormadapter.NewAdapter("sqlite3", "file::memory:?mode=memory&cache=shared"))
	e.EnableLog(true)
	e.EnableAutoSave(true)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e.AddPolicy(tt.userID.String(), fmt.Sprintf("/%s", tt.userID), http.MethodGet)

			if tt.disableACL {
				e.EnableEnforce(false)
			}

			r := BuildRouter(tt.fields.users, ja, e)
			token, _, _ := ja.Encode(jwt.MapClaims{"user_id": tt.userID})

			tt.args.r.Header.Set("Authorization", fmt.Sprintf("BEARER %s", token.Raw))
			r.ServeHTTP(tt.args.w, tt.args.r)

			w := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, tt.code, w.Code)
		})
	}
}
