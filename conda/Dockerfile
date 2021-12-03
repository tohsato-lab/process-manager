FROM continuumio/anaconda3:latest as backend-base

COPY docker/conda_packages/ /conda_packages/

RUN for package in `ls /conda_packages`; do conda env create -f="/conda_packages/$package"; done

RUN apt-get update -y && \
    apt-get install -y unzip ssh rsync sshpass && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

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

FROM backend-base

ARG UID=1000
ARG USERNAME=docker
RUN useradd -m -u ${UID} $USERNAME
RUN mkdir -p /process-manager/backend/src/ && \
    chown  $USERNAME:$USERNAME /process-manager/backend/src/

USER $USERNAME

WORKDIR /process-manager/backend/src
COPY --from=builder /go/bin/main ./main
COPY ./src/scripts/ ./scripts/

CMD ["./main"]
