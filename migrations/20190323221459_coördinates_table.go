package migration

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20190323221459, Down20190323221459)
}

func Up20190323221459(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(
`
create table coordinates (
  id varchar(50) not null primary key default uuid(),
  user_id varchar(50) not null,
  point point not null,
  foreign key user_id_fx (user_id) references users(id) on delete cascade
)
`)

	return errors.Wrap(err, "could not apply up migration to create coordinates table.")
}

func Down20190323221459(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`drop table coordinates`)

	return errors.Wrap(err, "could not apply down migration to drop coordinates table.")
}
