# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o watchtower ./cmd/watchtower

FROM alpine:3.18
RUN adduser -D watchtower
USER watchtower
COPY --from=build /src/watchtower /usr/local/bin/watchtower
ENTRYPOINT ["/usr/local/bin/watchtower"]
