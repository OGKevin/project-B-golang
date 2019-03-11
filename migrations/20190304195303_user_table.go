package migration

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
	"github.com/sirupsen/logrus"
)

func init() {
	goose.AddMigration(Up20190304195303, Down20190304195303)
}

func Up20190304195303(tx *sql.Tx) error {
	_, err := tx.Exec(`create table users (
  id varchar(50) not null default uuid() primary key,
  username varchar(50) not null unique,
  password varchar(64) not null
)`)
	if err != nil {
		logrus.WithError(err).Error("could not create user table")
		return errors.Wrap(err, "could not create user table")
	}

	return nil
}

func Down20190304195303(tx *sql.Tx) error {
	_, err := tx.Exec(`drop table users`)
	if err != nil {
		logrus.WithError(err).Error("could not drop table users")
		return errors.Wrap(err, "could not drop table users")
	}

	return nil
}
