//go:generate gojay -s=$GOFILE -t=NotFound,Error -o=generated_$GOFILE
package responses

import (
	"github.com/francoispqt/gojay"
	"github.com/sirupsen/logrus"
	"net/http"
)

type NotFound struct {
	// Ack Defines if the request was successful or not.
	Ack Ack `gojay:"ack,object"`
	// Error Explains why the server is responding with a bad request.
	Error *Error `gojay:"error"`
}

func newNotFound(ack Ack, error *Error) *NotFound {
	return &NotFound{Ack: ack, Error: error}
}

// WriteNotFound Returns a status code 404 with the NotFound Body
func WriteNotFound(w http.ResponseWriter, errr *Error) {
	errr.Code = http.StatusNotFound
	o := newNotFound(Ack{Ack: false}, errr)
	w.WriteHeader(http.StatusNotFound)
	err := gojay.NewEncoder(w).EncodeObject(o)
	if err != nil {
		logrus.WithError(err).Error("could not write NotFound response")
	}
}
