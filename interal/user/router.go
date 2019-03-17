package user

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

func BuildRouter(users Users, ja *jwtauth.JWTAuth) chi.Router {
	r := chi.NewRouter()

	r.Post("/", NewCreateUserRequest(users).ServeHttp)

	r.Route("/", func(r chi.Router) {
		r.Use(jwtauth.Verifier(ja))
		r.Use(jwtauth.Authenticator)

		r.Get("/{userId}", newGetUser(users).ServeHTTP)
	})

	return r
}
