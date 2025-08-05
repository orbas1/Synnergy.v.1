# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o synnergy ./cmd/synnergy

FROM alpine:3.18
RUN adduser -D synnergy
USER synnergy
COPY --from=build /src/synnergy /usr/local/bin/synnergy
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/synnergy"]
