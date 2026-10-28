FROM golang:1.27.3 AS go-builder

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

FROM node:26-bookworm-slim AS web-builder

WORKDIR /app

COPY web/package.json web/pnpm-lock.yaml web/pnpm-workspace.yaml ./web/
RUN npm install --global "$(node -p "require('./web/package.json').packageManager")" && (cd web && pnpm fetch)

COPY build.sh ./
COPY web/ ./web/
RUN ./build.sh --ui

FROM node:26-bookworm-slim AS runtime

RUN apt-get update && \
      apt-get install -y curl && \
      apt-get install -y --no-install-recommends \
      ca-certificates \
      awscli \
      jq \
      build-essential \
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