FROM nvidia/cuda:10.1-cudnn7-devel-ubuntu18.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y software-properties-common && \
    add-apt-repository ppa:longsleep/golang-backports && \
    apt-get update && \
    apt-get install -y less man curl wget git unzip bzip2 vim emacs && \
    apt-get install -y build-essential gcc clang zlib1g-dev libeigen3-dev cmake && \
    apt-get install -y python3 python3-dev python3-pip python3-venv python3-setuptools

# install anaconda
RUN wget -P /opt https://repo.anaconda.com/archive/Anaconda3-2020.02-Linux-x86_64.sh && \
    bash /opt/Anaconda3-2020.02-Linux-x86_64.sh -b -p /opt/anaconda3 && \
    rm /opt/Anaconda3-2020.02-Linux-x86_64.sh && \
    touch /root/.bashrc && \
    echo "export PATH=/opt/anaconda3/bin:$PATH" >> /root/.bashrc && \
    . /root/.bashrc && \
    conda init

# install server env
RUN apt-get install -y mysql-server golang-go

# install angular
RUN apt install -y nodejs npm && \
    npm install n -g && \
    n stable && \
    . /root/.bashrc && \
    echo n | npm install -g --silent @angular/cli

# tidying up
RUN apt-get clean && \
    apt-get autoclean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN git clone https://github.com/tohsato-lab/process-manager.git && \
    go get github.com/go-sql-driver/mysql && \
    go get github.com/gorilla/websocket
