package user

import (
	"context"
	"github.com/go-chi/jwtauth"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

func SetUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		if val, ok := claims["user_id"]; ok {
			r = r.WithContext(context.WithValue(r.Context(), "user_id", uuid.FromStringOrNil(val.(string))))
			next.ServeHTTP(w, r)
			return
		}

		logrus.Error("jwt token seems to not have an user id")
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}
