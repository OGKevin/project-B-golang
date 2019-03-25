package migration

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20190325091557, Down20190325091557)
}

func Up20190325091557(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
create unique index coordinates_user_id_point_uindex
	on coordinates (user_id, point);
`)
	return errors.Wrap(err, "could not add unique point index")
}

func Down20190325091557(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`
SET FOREIGN_KEY_CHECKS=false;
drop index coordinates_user_id_point_uindex on coordinates;
SET FOREIGN_KEY_CHECKS=true;`)
	
	return errors.Wrap(err, "could not revert unique point index creation")
}
