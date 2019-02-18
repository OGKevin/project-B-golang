FROM golang:1.11 as builder

WORKDIR /go/src/github.com/OGKevin/automation
COPY . .

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /http github.com/OGKevin/automation/cmd/http

FROM alpine

ARG TAG
ENV TAG=${TAG}
ENV SENTRY_RELEASE=${TAG}

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

WORKDIR /app

COPY --from=builder /http /app/http

CMD /bin/sh
