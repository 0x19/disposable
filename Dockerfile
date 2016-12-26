# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.7.1-alpine

ENV NAME=disposable
ENV DIR=/go/src/github.com/0x19/$NAME

RUN apk update && apk add --no-cache git

ADD . $DIR
WORKDIR $DIR

ENV GRPC_ADDR ":7290"
ENV HTTP_ADDR ":6290"

EXPOSE $GRPC_ADDR
EXPOSE $HTTP_ADDR

RUN go get
RUN go build && go install

ENTRYPOINT $NAME
