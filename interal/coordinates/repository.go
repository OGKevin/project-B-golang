package coordinates

import (
	"github.com/jmoiron/sqlx"
	"github.com/paulmach/go.geo"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type Point struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Coordinates *geo.Point
}

func NewPoint(userID uuid.UUID, coordinates *geo.Point) *Point {
	return &Point{UserID: userID, Coordinates: coordinates}
}

type coordinates interface {
	Create(point *Point) (uuid.UUID, error)
	Get(ID uuid.UUID) (Point, error)
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
	rows, err := c.db.Query(`select user_id, point from coordinates where id := ? limit 1`, ID)
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
	panic("implement me")
}

func (c *coordinatesFromDatabase) Delete(ID uuid.UUID) error {
	panic("implement me")
}
