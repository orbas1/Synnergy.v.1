# syntax=docker/dockerfile:1.4
FROM golang:1.21-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -o /opcodegen ./cmd/opcodegen

FROM alpine:3.18
RUN adduser -D -u 10001 opcodegen
USER opcodegen
WORKDIR /work
COPY --from=build /opcodegen /usr/local/bin/opcodegen
ENTRYPOINT ["/usr/local/bin/opcodegen"]
