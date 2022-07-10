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

FROM delfer/alpine-ftp-server

RUN DEBIAN_FRONTEND=noninteractive
RUN apk --no-cache add tzdata lftp bash \
    && rm -rf pidproxy
ENV TZ Asia/Tokyo

WORKDIR /home/process-manager/backend/src
COPY --from=builder /go/bin/main ./main
COPY ./src/scripts/ ./scripts/

CMD sh -c "/sbin/tini -- /bin/start_vsftpd.sh & ./main"
