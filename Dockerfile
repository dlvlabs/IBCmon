FROM golang:1.24.1-bookworm AS build

LABEL org.opencontainers.image.source="https://github.com/dlvlabs/IBCmon"

WORKDIR /build

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build

FROM debian:bookworm-slim AS deploy

RUN groupadd -r ibcmon && \
    useradd -r -g ibcmon ibcmon
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /build/ibcmon /usr/local/bin/ibcmon

USER ibcmon
ENTRYPOINT ["/usr/local/bin/ibcmon"]
