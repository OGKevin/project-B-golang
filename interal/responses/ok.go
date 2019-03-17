//go:generate gojay -s=$GOFILE -t=Created -o=generated_$GOFILE
//go:generate goimports -w  generated_$GOFILE
package responses

import (
	"github.com/francoispqt/gojay"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Created struct {
	// Ack Defines if the request was successful or not.
	Ack Ack       `gojay:"ack,object"`
	// ID The id of the newly created entity.
	ID  uuid.UUID `gojay:"id"`
}

// NewCreated creates a response object indicating that the entity got created
func NewCreated(ID uuid.UUID) *Created {
	return &Created{ID: ID, Ack:Ack{Ack: true}}
}

// WriteCreated writes created response indicating to the request
// that the entity has been created.
func WriteCreated(w http.ResponseWriter, ID uuid.UUID) {
	c := NewCreated(ID)
	w.WriteHeader(http.StatusCreated)
	err := gojay.NewEncoder(w).EncodeObject(c)
	if err != nil {
		logrus.WithError(err).Error("could not write created response")
	}
}