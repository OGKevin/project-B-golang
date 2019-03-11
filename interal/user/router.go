package user

import (
	"github.com/go-chi/chi"
)

func BuildRouter(users Users) chi.Router {
	r := chi.NewRouter()

	r.Post("/", NewCreateUserRequest(users).ServeHttp)

	return r
}
