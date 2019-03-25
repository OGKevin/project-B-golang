package coordinates

import (
	"encoding/json"
	"fmt"
	"github.com/OGKevin/project-B-golang/interal/responses"
	"github.com/asaskevich/govalidator"
	"github.com/casbin/casbin"
	"github.com/paulmach/go.geo"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type createBody struct {
	Longitude string `json:"longitude"valid:"longitude"`
	Latitude  string `json:"latitude"valid:"latitude"`
}

// create Creates new coordinates
// @Summary Save coordinates
// @Description Save coordinates
// @ID coordinates-create
// @Tags coordinates
// @Accept  json
// @Produce  json
// @Param body body coordinates.createBody true "The expected request body."
// @Success 200 {object} responses.Created "The response will include the id of the newly created user."
// @Failure 400 {object} responses.BadRequest "The error object will explain why the request failed."
// @Router /coordinates [post]
func create(coordinates coordinates, e *casbin.Enforcer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b createBody
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			logrus.WithError(err).Error("could not decode body")
			responses.WriteBadRequests(w, responses.NewError("Request body seems to not be a valid json."))
			return
		}

		if valid, err := govalidator.ValidateStruct(b); !valid {
			responses.WriteBadRequests(w, responses.NewValidationError(err.Error()))
			return
		}

		userID := r.Context().Value("user_id").(uuid.UUID)

		long, err := strconv.ParseFloat(b.Longitude, 64)
		if err != nil {
			logrus.WithError(err).Error("could not parse float")
			responses.WriteBadRequests(w, responses.NewErrorf("could not parse %q into float", b.Longitude))
		}

		lat, err := strconv.ParseFloat(b.Latitude, 64)
		if err != nil {
			logrus.WithError(err).Error("could not parse float")
			responses.WriteBadRequests(w, responses.NewErrorf("could not parse %q into float", b.Latitude))
		}

		ID, err := coordinates.Create(NewPoint(userID, geo.NewPoint(long, lat)))
		if err != nil {
			logrus.WithError(err).Error("saving coordinates failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		e.AddPolicy(userID.String(), fmt.Sprintf("*/coordinates/%s", ID), fmt.Sprintf("(%s)|(%s)", http.MethodGet, http.MethodDelete))
		e.AddPolicy(userID.String(), ".+/coordinates", fmt.Sprintf("(%s)", http.MethodGet))

		responses.WriteCreated(w, ID)
	}
}
