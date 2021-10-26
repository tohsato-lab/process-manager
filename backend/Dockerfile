FROM golang:latest as builder

WORKDIR /go/src

COPY ./src/go.mod /go/src/
RUN go mod download

COPY ./src/  /go/src/

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build \
    -o /go/bin/main \
    -ldflags '-s -w'

FROM continuumio/anaconda3

RUN apt update && \
    apt install unzip

ARG UID=1000
ARG USERNAME=docker
RUN useradd -m -u ${UID} $USERNAME
RUN mkdir -p /home/$USERNAME/process-manager/backend/src/ && \
    chown  $USERNAME:$USERNAME /home/$USERNAME/process-manager/backend/src/

USER $USERNAME

WORKDIR /home/$USERNAME/process-manager/backend/src
COPY --from=builder /go/bin/main ./main
COPY ./src/scripts/ ./scripts/

CMD ["./main"]
