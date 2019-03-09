package responses

import (
	"bytes"
	"github.com/francoispqt/gojay"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newResponseWriter() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

func TestWriteBadRequests(t *testing.T) {
	type args struct {
		w    http.ResponseWriter
		errr *Error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "main",
			args: args{
				w: newResponseWriter(),
				errr: NewError("error"),
			},
		},
		{
			name: "generated error message",
			args: args{
				w: newResponseWriter(),
				errr: NewError(uuid.NewV4().String()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteBadRequests(tt.args.w, tt.args.errr)

			recorder := tt.args.w.(*httptest.ResponseRecorder)
			if !assert.Equal(t, 400, recorder.Code) {
				return
			}

			buf := bytes.NewBufferString("")
			err := gojay.NewEncoder(buf).EncodeObject(newBadRequest(Ack{false}, tt.args.errr))
			if !assert.NoError(t, err) {
				return
			}

			if !assert.Equal(t, buf.String(), recorder.Body.String()) {
				return
			}
		})
	}
}
