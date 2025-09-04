# syntax=docker/dockerfile:1

# Build stage for Go binaries
ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build Synnergy node
RUN CGO_ENABLED=0 go build -o /out/synnergy ./cmd/synnergy
# Build wallet server
RUN CGO_ENABLED=0 go build -o /out/walletserver ./walletserver

# Runtime image for Synnergy node
FROM alpine:3.19 AS synnergy
COPY --from=build /out/synnergy /usr/local/bin/synnergy
EXPOSE 8080
ENTRYPOINT ["synnergy"]

# Runtime image for wallet server
FROM alpine:3.19 AS walletserver
COPY --from=build /out/walletserver /usr/local/bin/walletserver
EXPOSE 8090
ENTRYPOINT ["walletserver"]
