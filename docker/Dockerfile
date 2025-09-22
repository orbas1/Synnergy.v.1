# syntax=docker/dockerfile:1.5

ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-alpine AS build
RUN apk add --no-cache ca-certificates build-base git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go test ./...
RUN go build -trimpath -ldflags "-s -w" -o /out/synnergy ./cmd/synnergy
RUN go build -trimpath -ldflags "-s -w" -o /out/walletserver ./walletserver

FROM alpine:3.19 AS runtime-base
RUN addgroup -S synnergy && adduser -S -G synnergy synnergy
RUN apk add --no-cache ca-certificates tzdata wget
USER synnergy:synnergy
WORKDIR /var/lib/synnergy
ENV SYN_CONFIG=/etc/synnergy/config.yaml
COPY --from=build /out/synnergy /usr/local/bin/synnergy
COPY --from=build /out/walletserver /usr/local/bin/walletserver

FROM runtime-base AS synnergy
ENTRYPOINT ["/usr/local/bin/synnergy"]
EXPOSE 8080 30303 9102
HEALTHCHECK --interval=30s --timeout=5s --start-period=20s CMD wget -qO- http://127.0.0.1:8080/livez || exit 1

FROM runtime-base AS walletserver
ENTRYPOINT ["/usr/local/bin/walletserver"]
EXPOSE 8081 9103
HEALTHCHECK --interval=30s --timeout=5s --start-period=20s CMD wget -qO- http://127.0.0.1:8081/livez || exit 1
