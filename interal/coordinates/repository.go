package coordinates

import (
	"github.com/jmoiron/sqlx"
	"github.com/paulmach/go.geo"
	"github.com/satori/go.uuid"
)

type Point struct {
	ID uuid.UUID
	UserID uuid.UUID
	Coordinates *geo.Point
}

type coordinates interface {
	Create(point *geo.Point) (uuid.UUID, error)
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

func (coordinatesFromDatabase) Create(point *geo.Point) (uuid.UUID, error) {
	panic("implement me")
}

func (coordinatesFromDatabase) Get(ID uuid.UUID) (Point, error) {
	panic("implement me")
}

func (coordinatesFromDatabase) ListByUserID(userID uuid.UUID) (chan Point, error) {
	panic("implement me")
}

func (coordinatesFromDatabase) Delete(ID uuid.UUID) error {
	panic("implement me")
}
