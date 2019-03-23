package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OGKevin/project-B-golang/interal/acl"
	"github.com/OGKevin/xorm-adapter"
	"github.com/asaskevich/govalidator"
	"github.com/francoispqt/gojay"
	"github.com/go-chi/jwtauth"
	_ "github.com/mattn/go-sqlite3"
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
				r: httptest.NewRequest(http.MethodPost, "/", createUserBody(t, uuid.NewV4())),
			},
			code:  201,
			users: &usersForTest{},
		},
		{
			name: "username not unique",
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodPost, "/", createUserBody(t, uuid.NewV4())),
			},
			code:  400,
			users: &usersForTest{usernameNotUnique: true},
		},
		{
			name: "user creation failed",
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodPost, "/", createUserBody(t, uuid.NewV4())),
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
			t.Parallel()
			m := NewRouter(tt.users, ja, e)
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

func createUserBody(t *testing.T, pass uuid.UUID) io.Reader {
	buf := bytes.NewBufferString("")
	err := json.NewEncoder(buf).Encode(&userRequest{Username: uuid.NewV4().String(), Password: pass.String()})
	assert.NoError(t, err)

	return buf
}

type usersForTest struct {
	usernameNotUnique bool
	creationFailed    bool

	userNotFound bool

	password uuid.UUID
}

func (u *usersForTest) GetByUsername(username string) (*User, error) {
	var pass uuid.UUID

	if u.password == uuid.Nil {
		pass = uuid.NewV4()
	} else {
		pass = u.password
	}

	str := pass.String()

	p , err:= bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	err = bcrypt.CompareHashAndPassword(p, []byte(str))
	if err != nil {
		panic(err)
	}

	return &User{ID: uuid.NewV4(), Username: username, Password: string(p)}, nil
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
		name       string
		fields     fields
		args       args
		code       int
		userID     uuid.UUID
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
			code:   http.StatusOK,
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
			code:   http.StatusForbidden,
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
			code:       http.StatusBadRequest,
			userID:     uuid.NewV4(),
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
			code:       http.StatusNotFound,
			userID:     uuid.NewV4(),
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
			code:   http.StatusNotFound,
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

			r := NewRouter(tt.fields.users, ja, e)
			token, _, _ := ja.Encode(jwt.MapClaims{"user_id": tt.userID})

			tt.args.r.Header.Set("Authorization", fmt.Sprintf("BEARER %s", token.Raw))
			r.ServeHTTP(tt.args.w, tt.args.r)

			w := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, tt.code, w.Code)
		})
	}
}

func Test_login_ServeHTTP(t *testing.T) {
	type fields struct {
		users Users
		ja    *jwtauth.JWTAuth
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}

	ja := jwtauth.New("HS256", []byte("secret"), nil)

	tests := []struct {
		name   string
		fields fields
		args   args
		code int
	}{
		{
			name: "login ok",
			fields: fields{
				users: &usersForTest{password: uuid.FromStringOrNil("64EAE7E1-E192-434C-A314-C55FB3579C3A")},
				ja: ja,
			},
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodPost, "/login", createUserBody(t, uuid.FromStringOrNil("64EAE7E1-E192-434C-A314-C55FB3579C3A"))),
			},
			code: http.StatusOK,
		},
		{
			name: "login invalid password",
			fields: fields{
				users: &usersForTest{password: uuid.FromStringOrNil("6967776D-5444-417E-BE40-0F4C61DC7F89")},
				ja: ja,
			},
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodPost, "/login", createUserBody(t, uuid.FromStringOrNil("64EAE7E1-E192-434C-A314-C55FB3579C3A"))),
			},
			code: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter(tt.fields.users, tt.fields.ja, nil)

			r.ServeHTTP(tt.args.w, tt.args.r)

			w := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, tt.code, w.Code)
		})
	}
}
