package wellknown

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "main",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/.well-known/acme-challenge/34O4U9dXYQ95ZKzZagDAc34sigKhv4Verg6cPy1Yes8", nil),
			},
		},
	}

	assert.NoError(t, os.Setenv("ACME_TOKEN", "N9PAWbUUn_dXlM4v-ZTiW31_LK2O6xextoF6AUj5jFw"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := chi.NewMux()
			m.Route("/.well-known", Router)
			m.ServeHTTP(tt.args.w, tt.args.r)

			r := tt.args.w.(*httptest.ResponseRecorder)

			if !assert.Equal(t, http.StatusOK, r.Code) {
				return
			}

			assert.Equal(t, fmt.Sprintf("%s.%s", "34O4U9dXYQ95ZKzZagDAc34sigKhv4Verg6cPy1Yes8", "N9PAWbUUn_dXlM4v-ZTiW31_LK2O6xextoF6AUj5jFw"), r.Body.String())
		})
	}
}
