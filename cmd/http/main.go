package main

import (
	_ "github.com/OGKevin/project-B-golang/docs"
	"github.com/OGKevin/project-B-golang/interal/acl"
	"github.com/OGKevin/project-B-golang/interal/coordinates"
	"github.com/OGKevin/project-B-golang/interal/database"
	"github.com/OGKevin/project-B-golang/interal/logging"
	"github.com/OGKevin/project-B-golang/interal/response"
	"github.com/OGKevin/project-B-golang/interal/user"
	"github.com/OGKevin/project-B-golang/interal/wellknown"
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/docgen"
	"github.com/go-chi/jwtauth"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/OGKevin/xorm-adapter"
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
		if strings.Contains(r.RequestURI, "/api") || strings.Contains(r.RequestURI, ".json") {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		} else if strings.Contains(r.RequestURI, ".html") {
			w.Header().Set("Content-Type", "text/html")
			next.ServeHTTP(w, r)
		} else if strings.Contains(r.RequestURI, ".js") {
			w.Header().Set("Content-Type", "text/javascript")
			next.ServeHTTP(w, r)
		} else if strings.Contains(r.RequestURI, ".css") {
			w.Header().Set("Content-Type", "text/css")
			next.ServeHTTP(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		}
	})
}

func main() {
	logging.Trace(logging.TraceTypeEntering)
	defer logging.Trace(logging.TraceTypeExiting)

	db := database.GetDB()

	r := createRouter(db)

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

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @host project-b.ogkevin.nl
// @BasePath /api/v1
func createRouter(db *sqlx.DB) *chi.Mux{
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(setContentTypeHeader)
	r.Use(middleware.Timeout(time.Second * 5))

	r.Get("/", index)
	r.Mount("/api/v1", apiRouter(db))
	r.Route("/.well-known", wellknown.Router)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	docgen.PrintRoutes(r)

	return r
}

func apiRouter(db *sqlx.DB) chi.Router {
	ja := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)
	e := acl.NewEnforcer(xormadapter.NewAdapter("mysql", os.Getenv("DB_PATH"), true))
	e.EnableAutoSave(true)

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/doc.json", http.StatusPermanentRedirect)
	})
	r.Mount("/user", user.NewRouter(user.NewUsersDatabase(db), ja, e))
	r.Mount("/coordinates", coordinates.NewRouter(coordinates.NewCoordinatesFromDatabase(db), ja, e))
	return r
}

func index(w http.ResponseWriter, _ *http.Request) {
	err := response.WriteAckTrue(w)
	if err != nil {
		logrus.WithError(err).Error("could not send index")
	}
}
