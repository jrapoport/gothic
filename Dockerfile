FROM golang:1.15-alpine as build
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apk add --no-cache make git protobuf

WORKDIR /go/src/github.com/jrapoport/gothic

# Pulling dependencies
COPY ./Makefile ./go.* ./
RUN make deps

# Building stuff
COPY . /go/src/github.com/jrapoport/gothic
RUN make build

FROM alpine:3.7
RUN adduser -D -u 1000 gothic

RUN apk add --no-cache ca-certificates
COPY --from=build /go/src/github.com/jrapoport/gothic/gothic /usr/local/bin/gothic

USER gothic
CMD ["gothic"]
