FROM golang:1.24.6 AS go-builder

RUN apt-get update && apt-get install -y --no-install-recommends \
      build-essential \
      binutils \
      binutils-gold \
      pkg-config \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/github.com/patterninc/heimdall

COPY go.mod go.sum ./
RUN go mod download

COPY build.sh ./
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY plugins ./plugins

# Declared right before its only consumer so it invalidates as little as possible.
ARG BUILD_VERSION=dev
RUN BUILD_VERSION="${BUILD_VERSION}" ./build.sh --go

FROM node:20-bookworm AS web-builder

WORKDIR /app

COPY web/pnpm-lock.yaml web/pnpm-workspace.yaml ./web/
RUN corepack enable && corepack prepare pnpm@10.28.2 --activate && (cd web && pnpm fetch)

COPY build.sh ./
COPY web/ ./web/
RUN ./build.sh --ui

FROM node:20-bookworm-slim AS runtime

RUN apt-get update && \
      apt-get install -y curl && \
      apt-get install -y --no-install-recommends \
      ca-certificates \
      awscli \
      jq \
      build-essential \
 && corepack enable && corepack prepare pnpm@10.28.2 --activate \
 && rm -rf /var/lib/apt/lists/*

COPY --from=go-builder /usr/local/go /usr/local/go
ENV PATH="/usr/local/go/bin:${PATH}"
ENV CGO_ENABLED=1

WORKDIR /go/src/github.com/patterninc/heimdall

COPY --from=go-builder /go/src/github.com/patterninc/heimdall .

# Rarely-changing runtime files first, frequently-changing build outputs last.
COPY configs/local.yaml /etc/heimdall/heimdall.yaml
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
COPY assets ./assets
COPY --from=go-builder /go/src/github.com/patterninc/heimdall/dist ./dist
COPY --from=web-builder /app/web ./web

CMD [ "/usr/local/bin/entrypoint.sh" ]