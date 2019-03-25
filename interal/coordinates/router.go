package coordinates

import (
	"github.com/OGKevin/project-B-golang/interal/acl"
	"github.com/OGKevin/project-B-golang/interal/user"
	"github.com/casbin/casbin"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

func NewRouter(coordinates coordinates, ja *jwtauth.JWTAuth, e *casbin.Enforcer) *chi.Mux {
	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(ja))
	r.Use(jwtauth.Authenticator)
	r.Use(user.SetUserID)

	r.Post("/", create(coordinates, e))
	r.Route("/", func(r chi.Router) {
		r.Use(acl.BuildPermissionCheckMiddleware(e))
		r.Get("/{id}", get(coordinates))
	})

	return r
}
