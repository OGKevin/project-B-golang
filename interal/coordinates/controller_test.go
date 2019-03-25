package coordinates

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/paulmach/go.geo"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OGKevin/project-B-golang/interal/acl"
	"github.com/OGKevin/xorm-adapter"
	"github.com/casbin/casbin"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_create(t *testing.T) {
	type args struct {
		coordinates coordinates
		e           *casbin.Enforcer
	}
	tests := []struct {
		name        string
		args        args
		response    http.ResponseWriter
		request     *http.Request
		code        int
		coordinates coordinates
	}{
		{
			name:        "ok created",
			code:        http.StatusCreated,
			response:    httptest.NewRecorder(),
			request:     httptest.NewRequest(http.MethodPost, "/", createRequestBody(t)),
			coordinates: &coordinatesForTest{},
		},
		{
			name:        "no body",
			code:        http.StatusBadRequest,
			response:    httptest.NewRecorder(),
			request:     httptest.NewRequest(http.MethodPost, "/", nil),
			coordinates: &coordinatesForTest{},
		},
		{
			name:        "bad body",
			code:        http.StatusBadRequest,
			response:    httptest.NewRecorder(),
			request:     httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{\"latitude\": \"dhgjasdhjaks\"}")),
			coordinates: &coordinatesForTest{},
		},
	}

	ja := jwtauth.New("HS256", []byte("secret"), nil)
	e := acl.NewEnforcer(xormadapter.NewAdapter("sqlite3", "file::memory:?mode=memory&cache=shared"))
	e.EnableLog(true)
	e.EnableAutoSave(true)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, _, err := ja.Encode(jwt.MapClaims{"user_id": uuid.NewV4()})
			if !assert.NoError(t, err) {
				return
			}

			tt.request.Header.Set("Authorization", fmt.Sprintf("BEARER %s", token.Raw))
			r := NewRouter(tt.coordinates, ja, e)
			r.ServeHTTP(tt.response, tt.request)

			w := tt.response.(*httptest.ResponseRecorder)
			if !assert.Equal(t, tt.code, w.Code, w.Body.String()) {
				return
			}
		})
	}
}

type coordinatesForTest struct {
	getNotFound bool
}

func (*coordinatesForTest) Create(point *Point) (uuid.UUID, error) {
	return uuid.NewV4(), nil
}

func (c *coordinatesForTest) Get(ID uuid.UUID) (*Point, error) {
	if c.getNotFound  {
		return nil, errors.New("not found")
	}

	return NewPoint(uuid.NewV4(), geo.NewPoint(10.000, 50.000)), nil
}

func (*coordinatesForTest) ListByUserID(userID uuid.UUID) (chan Point, error) {
	panic("implement me")
}

func (*coordinatesForTest) Delete(ID uuid.UUID) error {
	panic("implement me")
}

func createRequestBody(t *testing.T) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(createBody{Longitude: "4.895168", Latitude: "52.370216"})
	if !assert.NoError(t, err) {
		return nil
	}

	return buf

}

func Test_get(t *testing.T) {
	type args struct {
		coordinates coordinates
	}
	tests := []struct {
		name string
		args args
		response    http.ResponseWriter
		request     *http.Request
		code        int
		disableEnforcer bool
	}{
		{
			name: "get non existing",
			args: args{
				coordinates: &coordinatesForTest{getNotFound: true},
			},
			response: httptest.NewRecorder(),
			request: httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", uuid.NewV4()), nil),
			code: http.StatusNotFound,
			disableEnforcer: true,
		},
		{
			name: "get existing",
			args: args{
				coordinates: &coordinatesForTest{},
			},
			response: httptest.NewRecorder(),
			request: httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", uuid.NewV4()), nil),
			code: http.StatusOK,
			disableEnforcer: true,
		},
	}

	ja := jwtauth.New("HS256", []byte("secret"), nil)
	e := acl.NewEnforcer(xormadapter.NewAdapter("sqlite3", "file::memory:?mode=memory&cache=shared"))
	e.EnableLog(true)
	e.EnableAutoSave(true)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.disableEnforcer {
				e.EnableEnforce(false)
			}

			r := NewRouter(tt.args.coordinates, ja, e)

			token, _, err := ja.Encode(jwt.MapClaims{"user_id": uuid.NewV4()})
			if !assert.NoError(t, err) {
				return
			}

			tt.request.Header.Set("Authorization", fmt.Sprintf("BEARER %s", token.Raw))

			r.ServeHTTP(tt.response, tt.request)

			w := tt.response.(*httptest.ResponseRecorder)

			if !assert.Equal(t, tt.code, w.Code) {
				return
			}
		})
	}
}
