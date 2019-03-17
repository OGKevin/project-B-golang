package user

import (
	"bytes"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/francoispqt/gojay"
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := BuildRouter(tt.users, jwtauth.New("HS256", []byte("secret"), nil))
			m.ServeHTTP(tt.args.w, tt.args.r)

			w := tt.args.w.(*httptest.ResponseRecorder)
			if !assert.Equal(t, tt.code, w.Code) {
				return
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
	}{
		{
			name: "get user",
			fields: fields{
				users: &usersForTest{},
			},
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", uuid.NewV4()), nil),
			},
			code: http.StatusOK,
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
			code: http.StatusBadRequest,
		},
		{
			name: "get user with valid id not found",
			fields: fields{
				users: &usersForTest{userNotFound: true},
			},
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", uuid.NewV4()), nil),
			},
			code: http.StatusNotFound,
		},
	}

	ja := jwtauth.New("HS256", []byte("secret"), nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := BuildRouter(tt.fields.users, ja)
			token, _, _ := ja.Encode(jwt.MapClaims{"some": "user"})

			tt.args.r.Header.Set("Authorization", fmt.Sprintf("BEARER %s", token.Raw))
			r.ServeHTTP(tt.args.w, tt.args.r)

			w := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, tt.code, w.Code)
		})
	}
}
