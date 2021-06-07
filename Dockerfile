FROM golang:alpine as build-stage
ENV GO111MODULE=on
ENV CGO_ENABLED=0

RUN apk add --no-cache make git

WORKDIR /gothic

# Pulling dependencies
COPY . .
RUN make deps

# Building stuff
RUN make release

FROM alpine:3.7
RUN adduser -D -u 1000 gothic

RUN apk add --no-cache ca-certificates
COPY --from=build-stage /gothic/build/release/gothic /usr/local/bin/gothic

USER gothic
CMD ["gothic"]
