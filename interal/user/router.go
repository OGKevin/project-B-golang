package user

import (
	"github.com/OGKevin/project-B-golang/interal/database"
	"github.com/go-chi/chi"
)

func Router(r chi.Router) {
	db := database.GetDB()
	r.Route("/user", func(r chi.Router) {
		r.Post("/", NewCreateUserRequest(NewUsersDatabase(db)).ServeHttp)
	})
}
