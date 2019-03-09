package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/OGKevin/project-B-golang/interal/logging"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var defaultDB *sqlx.DB
var once sync.Once

func ini() {
	logging.Trace(logging.TraceTypeEntering)
	defer logging.Trace(logging.TraceTypeExiting)

	db, err := sqlx.Open("mysql", os.Getenv("DB_PATH"))
	if err != nil {
		logrus.WithError(err).Error("could not open database connection")
	}

	defaultDB = db
}

// GetDB returns the default database connection
func GetDB() *sqlx.DB {
	once.Do(ini)
	logging.Trace(logging.TraceTypeEntering)
	defer logging.Trace(logging.TraceTypeExiting)

	return defaultDB
}
