# Prod Dockerfile with optimizations enabled
# Compile stage
ARG GO_VERSION=1.14

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

# Copy sources from current dir to Container
COPY . .

RUN go build -o application main.go

FROM alpine:3.12

WORKDIR /
COPY --from=builder /application .
ADD assets /assets

EXPOSE 8080

ENV GIN_MODE=release

ENTRYPOINT ["./application"]
