package main

import (
	"github.com/OGKevin/project-B-golang/interal/database"
	"github.com/OGKevin/project-B-golang/interal/logging"
	"github.com/OGKevin/project-B-golang/interal/response"
	"github.com/OGKevin/project-B-golang/interal/user"
	"github.com/asaskevich/govalidator"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/http-swagger"
	"net/http"
	_ "github.com/OGKevin/project-B-golang/docs"
	"strings"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Warning("could not load dot env")
	}

	govalidator.SetFieldsRequiredByDefault(true)
}

type setContentTypeHeader struct {
	handler http.Handler
}

func (s setContentTypeHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.RequestURI, "/api") || strings.Contains(r.RequestURI, ".json"){
		w.Header().Add("Content-Type", "application/json")
		s.handler.ServeHTTP(w, r)
	} else if strings.Contains(r.RequestURI, ".html") {
		w.Header().Add("Content-Type", "text/html")
		s.handler.ServeHTTP(w, r)
	}
}

func main() {
	logging.Trace(logging.TraceTypeEntering)
	defer logging.Trace(logging.TraceTypeExiting)

	router := createRouter(database.GetDB())

	err := http.ListenAndServe(":80", setContentTypeHeader{handler: router})
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
func createRouter(db *sqlx.DB) *httprouter.Router {
	router := httprouter.New()
	router.GET("/", index)
	router.GET("/api/v1", index)
	router.POST("/api/v1/user", user.NewCreateUserRequest(user.NewUsersDatabase(db)).Handle)
	router.GET("/swagger/*swagger", swagger)

	return router
}

func swagger(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	httpSwagger.WrapHandler.ServeHTTP(w, r)
}

func index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	err := response.WriteAckTrue(w)
	if err != nil {
		logrus.WithError(err).Error("could not send index")
	}
}
