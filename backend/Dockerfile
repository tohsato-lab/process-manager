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

FROM ubuntu:20.04

RUN DEBIAN_FRONTEND=noninteractive
RUN apt-get update && \
    apt-get install -y tzdata wget && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
ENV TZ Asia/Tokyo

WORKDIR /process-manager/backend/src
COPY --from=builder /go/bin/main ./main
COPY ./src/scripts/ ./scripts/

CMD ["./main"]
