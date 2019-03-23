package coordinates

import (
	"github.com/casbin/casbin"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"net/http"
)

func NewRouter(coordinates coordinates, ja *jwtauth.JWTAuth, e *casbin.Enforcer) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	return r
}
