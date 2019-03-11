package user

import (
	"bytes"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/francoispqt/gojay"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func newResponseWriter() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

func TestCreateUser(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		r   *http.Request
	}
	tests := []struct {
		name string
		args args
		code int
		users Users
	}{
		{
			name: "created",
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodPost, "/", createUserBody(t)),
			},
			code: 201,
			users: &usersForTest{},
		},
		{
			name: "username not unique",
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodPost, "/", createUserBody(t)),
			},
			code: 400,
			users: &usersForTest{usernameNotUnique: true},
		},
		{
			name: "user creation failed",
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodPost, "/", createUserBody(t)),
			},
			code: 500,
			users: &usersForTest{creationFailed: true},
		},
		{
			name: "no body",
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("")),
			},
			code: 400,
			users: &usersForTest{},
		},
		{
			name: "bad body",
			args: args{
				w: newResponseWriter(),
				r: httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{}")),
			},
			code: 400,
			users: &usersForTest{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := BuildRouter(tt.users)
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
	creationFailed bool
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
