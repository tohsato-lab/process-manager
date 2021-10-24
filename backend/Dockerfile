FROM continuumio/anaconda3

ARG USERNAME=docker
ARG GROUPNAME=docker
ARG UID=1000
ARG GID=1000
RUN groupadd -g $GID $GROUPNAME && \
    useradd -m -s /bin/bash -u $UID -g $GID $USERNAME

RUN apt update && \
    apt install unzip

USER $USERNAME
WORKDIR /home/$USERNAME/