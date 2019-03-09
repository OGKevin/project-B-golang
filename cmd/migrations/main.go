package main

import (
	"github.com/OGKevin/project-B-golang/interal/database"
	"github.com/OGKevin/project-B-golang/interal/logging"
	_ "github.com/OGKevin/project-B-golang/migrations"
	"github.com/joho/godotenv"
	"github.com/pressly/goose"
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Warning("could not load dot env")
	}

	err = goose.SetDialect("mysql")
	if err != nil {
		logrus.WithError(err).Fatal("could not set goose dialect")
	}
}

func main() {
	logging.Trace(logging.TraceTypeEntering)
	defer logging.Trace(logging.TraceTypeExiting)

	err := goose.Run(os.Args[1], database.GetDB().DB, "./migrations")
	if err != nil {
		logrus.WithError(err).Error("could not apply migrations")
		return
	}
}
