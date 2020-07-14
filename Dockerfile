ARG GO_VERSION=1.14
ARG APP_NAME=go-rest

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /api
WORKDIR /api

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /${APP_NAME} ./...

FROM apline:latest

RUN mkdir -p /api
WORKDIR /api
COPY --from=builder /api/${APP_NAME} .

EXPOSE 8080

ENTRYPOINT ["./${APP_NAME}"]
