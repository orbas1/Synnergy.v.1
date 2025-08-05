# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o opcodegen ./cmd/opcodegen

FROM alpine:3.18
RUN adduser -D opcodegen
USER opcodegen
WORKDIR /work
COPY --from=build /src/opcodegen /usr/local/bin/opcodegen
ENTRYPOINT ["/usr/local/bin/opcodegen"]
