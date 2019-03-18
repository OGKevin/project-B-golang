package acl

import (
	"github.com/casbin/casbin"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

func BuildPermissionCheckMiddleware(e *casbin.Enforcer) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := e.LoadPolicy()
			if err != nil {
				logrus.WithError(err).Error("could not load acl policies")
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			userId := r.Context().Value("user_id").(uuid.UUID)

			if !e.Enforce(userId.String(), r.URL.Path, r.Method) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
