# Dev Dockerfile with Debugging enabled and no optimizations
# Compile stage
ARG GO_VERSION=1.14

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download
# Copy sources from current dir to Container
COPY . .

# The -gcflags "all=-N -l" flag helps us get a better debug experience
RUN go build -gcflags "all=-N -l" -o application main.go

# Compile Delve (Go debugger)
RUN go get github.com/go-delve/delve/cmd/dlv

# Final stage
FROM alpine:3.12

# Allow delve to run on Alpine based containers.
RUN apk add --no-cache libc6-compat

WORKDIR /
COPY --from=builder /application .
COPY --from=builder /go/bin/dlv /
ADD assets /assets

# Port 8080 belongs to our application, 40000 belongs to Delve
EXPOSE 8080 40000

ENV GIN_MODE=debug

# Run delve and application
CMD ["/dlv", "--continue", "--accept-multiclient", "--listen=:40000", "--headless=true", "--api-version=2", "exec", "/application"]