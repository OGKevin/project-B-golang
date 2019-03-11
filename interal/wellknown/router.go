package wellknown

import (
	"github.com/go-chi/chi"
)

func Router(r chi.Router) {
	r.Get("/acme-challenge/{token}", serveHTTP)
}
