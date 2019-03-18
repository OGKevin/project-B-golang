package user

import (
	"github.com/OGKevin/project-B-golang/interal/acl"
	"github.com/casbin/casbin"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

func BuildRouter(users Users, ja *jwtauth.JWTAuth, e *casbin.Enforcer) chi.Router {
	r := chi.NewRouter()

	r.Post("/", NewCreateUserRequest(users, e).ServeHTTP)
	r.Post("/login", newLogin(users, ja).ServeHTTP)

	r.Route("/", func(r chi.Router) {
		r.Use(jwtauth.Verifier(ja))
		r.Use(jwtauth.Authenticator)
		r.Use(SetUserID)
		r.Use(acl.BuildPermissionCheckMiddleware(e))

		r.Get("/{userId}", newGetUser(users).ServeHTTP)
	})

	return r
}
