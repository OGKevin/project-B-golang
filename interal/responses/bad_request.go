//go:generate gojay -s=$GOFILE -t=BadRequest,Error -o=generated_$GOFILE
package responses

import (
	"fmt"
	"github.com/francoispqt/gojay"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Error struct {
	Message string `gojay:"message"`
	Code int `gojay:"code"`
}

// NewError Constructs a new error object
func NewError(message string) *Error {
	return &Error{Message: message}
}
// NewErrorf Constructs a new error object
func NewErrorf(message string, args ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(message, args...)}
}

// NewValidationError Constructs a new error object and appends
// the message with validation error.
func NewValidationError(message string) *Error {
	return &Error{Message: fmt.Sprintf("Validation error: %s", message)}
}

type BadRequest struct {
	Ack   Ack    `gojay:"ack,object"`
	Error *Error `gojay:"error"`
}

func newBadRequest(ack Ack, error *Error) *BadRequest {
	return &BadRequest{Ack: ack, Error: error}
}

// WriteBadRequests Returns a status code 400 with the BadRequest Body
func WriteBadRequests(w http.ResponseWriter, errr *Error) {
	errr.Code = http.StatusBadRequest
	o := newBadRequest(Ack{Ack:false}, errr)
	w.WriteHeader(http.StatusBadRequest)
	err := gojay.NewEncoder(w).EncodeObject(o)
	if err != nil {
		logrus.WithError(err).Error("could not write bad request response")
	}
}
