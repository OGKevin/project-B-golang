package wellknown

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	verification, er := getAcmeVerificationToken()

	if er != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := errors.Wrap(er, "wellknown: could not get Acme verification token")
		logrus.WithError(err).Error("error handling acme")
	} else {
		_, er := io.WriteString(w, fmt.Sprintf("%s.%s", chi.URLParam(r, "token"), verification))
		if er != nil {
			err := errors.Wrap(er, "wellknown: could not write response")
			logrus.WithError(err).Error("could not write response")
		}
	}
}

func getAcmeVerificationToken() (string, error) {
	token := os.Getenv("ACME_TOKEN")

	if token == "" {
		return "", errors.New("wellKnown: the env 'ACME-TOKEN' seems to be not set")
	}

	return token, nil
}
