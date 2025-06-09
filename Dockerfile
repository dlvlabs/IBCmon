FROM golang:1.24.1-alpine AS builder

WORKDIR /code/
COPY . /code/

ARG GOARCH
ARG GOOS=linux

RUN GOOS=${GOOS} GOARCH=${GOARCH} go build

FROM alpine:latest

COPY --from=builder /code/ibcmon /usr/bin/

ENTRYPOINT ["ibcmon"]
