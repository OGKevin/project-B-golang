FROM golang:1.11 as builder

WORKDIR /go/src/github.com/OGKevin/project-B-golang

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /http github.com/OGKevin/project-B-golang/cmd/http

FROM alpine

ARG TAG
ENV TAG=${TAG}
ENV SENTRY_RELEASE=${TAG}
ENV GOPATH=/go

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

RUN mkdir -p /go/src/github.com/OGKevin/project-B-golang
COPY --from=builder /go/src/github.com/OGKevin/project-B-golang /go/src/github.com/OGKevin/project-B-golang

WORKDIR /app

COPY --from=builder /http /http

CMD /bin/sh
