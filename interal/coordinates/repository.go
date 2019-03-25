package coordinates

import (
	"github.com/jmoiron/sqlx"
	"github.com/paulmach/go.geo"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type Point struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
 	Coordinates *geo.Point `json:"coordinates"`
}

func NewPoint(userID uuid.UUID, coordinates *geo.Point) *Point {
	return &Point{UserID: userID, Coordinates: coordinates}
}

type coordinates interface {
	Create(point *Point) (uuid.UUID, error)
	Get(ID uuid.UUID) (*Point, error)
	ListByUserID(userID uuid.UUID) (chan Point, error)
	Delete(ID uuid.UUID) error
}

type coordinatesFromDatabase struct {
	db *sqlx.DB
}

func NewCoordinatesFromDatabase(db *sqlx.DB) *coordinatesFromDatabase {
	return &coordinatesFromDatabase{db: db}
}

func (c *coordinatesFromDatabase) Create(point *Point) (uuid.UUID, error) {
	id := uuid.NewV4()
	_, err := c.db.Exec(`insert into coordinates (id, user_id, point) value (?, ?, geometryfromtext(?))`, id, point.UserID, point.Coordinates.ToWKT())
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "could not save coordinates in database")
	}

	return id, nil
}

func (c *coordinatesFromDatabase) Get(ID uuid.UUID) (*Point, error) {
	rows, err := c.db.Query(`select user_id, point from coordinates where id = ? limit 1`, ID)
	if err != nil {
		return nil, errors.Wrap(err, "could not get point by id")
	}

	defer rows.Close()

	var p Point

	rows.Next()
	if err = rows.Scan(&p.UserID, &p.Coordinates); err != nil {
		return nil, errors.Wrap(err, "could not scan result for getting point by id")
	}

	p.ID = ID

	return &p, nil
}

func (c *coordinatesFromDatabase) ListByUserID(userID uuid.UUID) (chan Point, error) {
	rows, err := c.db.Query(`select id, point from coordinates where user_id = ?`, userID)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch all coordinates by user")
	}

	defer rows.Close()

	ch := make(chan Point)

	go func() {
		defer close(ch)
		for rows.Next() {
			var p Point
			if err = rows.Scan(&p.ID, &p.Coordinates); err != nil {
				logrus.WithError(err).Errorf("could not scan rows for coordinates into struct")
				continue
			}
			p.UserID = userID

			ch <- p
		}
	}()

	return ch, nil
}

func (c *coordinatesFromDatabase) Delete(ID uuid.UUID) error {
	_, err := c.db.Exec(`delete from coordinates where id = ?`, ID)
	if err != nil {
		return errors.Wrap(err, "could not delete coordinate")
	}

	return nil
}
