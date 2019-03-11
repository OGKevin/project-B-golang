package main

import (
	_ "github.com/OGKevin/project-B-golang/docs"
	"github.com/OGKevin/project-B-golang/interal/logging"
	"github.com/OGKevin/project-B-golang/interal/response"
	"github.com/OGKevin/project-B-golang/interal/user"
	"github.com/OGKevin/project-B-golang/interal/wellknown"
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/http-swagger"
	"net/http"
	"strings"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Warning("could not load dot env")
	}

	govalidator.SetFieldsRequiredByDefault(true)
}

func setContentTypeHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.RequestURI, "/api") || strings.Contains(r.RequestURI, ".json"){
			w.Header().Add("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		} else if strings.Contains(r.RequestURI, ".html") {
			w.Header().Add("Content-Type", "text/html")
			next.ServeHTTP(w, r)
		} else if strings.Contains(r.RequestURI, ".js") {
			w.Header().Add("Content-Type", "text/javascript")
			next.ServeHTTP(w, r)
		} else if strings.Contains(r.RequestURI, ".css") {
			w.Header().Add("Content-Type", "text/css")
			next.ServeHTTP(w, r)
		}
	})
}

func main() {
	logging.Trace(logging.TraceTypeEntering)
	defer logging.Trace(logging.TraceTypeExiting)

	r := createRouter()

	err := http.ListenAndServe(":80", r)
	logrus.WithError(err).Fatal("http server died")
}

// @title Project b
// @version 1.0
// @description Well you know, nothing important. Just making sure people can capture memories
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host project-b.ogkevin.nl
// @BasePath /api/v1
func createRouter() *chi.Mux{
	r := chi.NewRouter()
	r.Use(setContentTypeHeader)

	r.Get("/", index)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", index)
		r.Route("/user", user.Router)
	})
	r.Get("/.well-known/acme-challenge/{token}", wellknown.ServeHTTP)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}

func index(w http.ResponseWriter, _ *http.Request) {
	err := response.WriteAckTrue(w)
	if err != nil {
		logrus.WithError(err).Error("could not send index")
	}
}
