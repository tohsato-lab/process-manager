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

FROM stilliard/pure-ftpd

RUN DEBIAN_FRONTEND=noninteractive
RUN apt-get update && \
    apt-get install -y tzdata lftp && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
ENV TZ Asia/Tokyo

WORKDIR /home/process-manager/backend/src
COPY --from=builder /go/bin/main ./main
COPY ./src/scripts/ ./scripts/

CMD sh -c "/run.sh -c 30 -C 10 -l puredb:/etc/pure-ftpd/pureftpd.pdb -E -j -R -P $PUBLICHOST -p 30000:30059 & ./main"
