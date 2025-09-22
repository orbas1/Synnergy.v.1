# syntax=docker/dockerfile:1

ARG GO_VERSION=1.22
ARG NODE_VERSION=20

FROM golang:${GO_VERSION}-alpine AS go-build
WORKDIR /src
RUN apk add --no-cache build-base git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/synnergy ./cmd/synnergy
RUN CGO_ENABLED=0 go build -o /out/walletserver ./walletserver

FROM node:${NODE_VERSION}-alpine AS web-build
WORKDIR /app
COPY web/package*.json ./
RUN npm ci --legacy-peer-deps
COPY web .
RUN npm run build

FROM alpine:3.19 AS runtime-base
RUN apk add --no-cache ca-certificates tzdata bash && \
    addgroup -S synnergy && adduser -S synnergy -G synnergy
ENV SYN_CONFIG=/etc/synnergy/synnergy.yaml
WORKDIR /home/synnergy
USER synnergy

FROM runtime-base AS synnergy
COPY --from=go-build /out/synnergy /usr/local/bin/synnergy
COPY configs/dev.yaml /etc/synnergy/synnergy.yaml
EXPOSE 8080 30303
HEALTHCHECK --interval=30s --timeout=5s CMD wget -qO- http://127.0.0.1:8080/healthz || exit 1
ENTRYPOINT ["/usr/local/bin/synnergy"]

FROM runtime-base AS walletserver
COPY --from=go-build /out/walletserver /usr/local/bin/walletserver
COPY walletserver/config.example.yaml /etc/synnergy/wallet.yaml
ENV SYN_WALLET_CONFIG=/etc/synnergy/wallet.yaml
EXPOSE 8090
HEALTHCHECK --interval=30s --timeout=5s CMD wget -qO- http://127.0.0.1:8090/healthz || exit 1
ENTRYPOINT ["/usr/local/bin/walletserver"]

FROM node:${NODE_VERSION}-alpine AS web
WORKDIR /app
ENV NODE_ENV=production
COPY --from=web-build /app/package*.json ./
COPY --from=web-build /app/.next ./.next
COPY --from=web-build /app/public ./public
RUN npm ci --omit=dev --legacy-peer-deps
EXPOSE 3000
CMD ["npm", "run", "start"]
