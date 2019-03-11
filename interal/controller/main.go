package controller

import (
	"net/http"
)

type Controller interface {
	ServeHttp(w http.ResponseWriter, r *http.Request)
}
