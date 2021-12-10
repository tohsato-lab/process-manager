FROM golang:1.17.2 as builder

WORKDIR /go/src

COPY ./src/go.* /go/src/
RUN go mod download

COPY ./src/  /go/src/

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build \
    -o /go/bin/main \
    -ldflags '-s -w'

FROM linuxserver/openssh-server

RUN apk add --no-cache rsync

COPY --from=builder /go/bin/main ./main

CMD ["./main"]
