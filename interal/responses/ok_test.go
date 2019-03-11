package responses

import (
	"bytes"
	"github.com/francoispqt/gojay"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestWriteCreated(t *testing.T) {
	type args struct {
		w  http.ResponseWriter
		ID uuid.UUID
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name:"main",
			args: args{
				w: newResponseWriter(),
				ID: uuid.NewV4(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteCreated(tt.args.w, tt.args.ID)

			recorder := tt.args.w.(*httptest.ResponseRecorder)
			if !assert.Equal(t, 201, recorder.Code) {
				return
			}

			buf := bytes.NewBufferString("")
			err := gojay.NewEncoder(buf).EncodeObject(NewCreated(tt.args.ID))
			if !assert.NoError(t, err) {
				return
			}

			if !assert.Equal(t, buf.String(), recorder.Body.String()) {
				return
			}
		})
	}
}
