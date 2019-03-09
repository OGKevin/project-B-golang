package controller

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Controller interface {
	Handle(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}
