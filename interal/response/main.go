package response

import (
	"encoding/json"
	"github.com/OGKevin/project-B-golang/interal/logging"
	"github.com/pkg/errors"
	"net/http"
)

type Ack struct {
	Ack bool `json:"ack"`
}

func WriteAckTrue(w http.ResponseWriter) error {
	logging.Trace(logging.TraceTypeEntering)
	defer logging.Trace(logging.TraceTypeExiting)

	err := json.NewEncoder(w).Encode(Ack{Ack: true})
	if err != nil {
		return errors.Wrap(err, "could not write ack true response")
	}

	return nil
}