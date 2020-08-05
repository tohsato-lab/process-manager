FROM nvidia/cuda:10.1-cudnn7-devel-ubuntu18.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y software-properties-common && \
    add-apt-repository ppa:longsleep/golang-backports && \
    apt-get update && \
    apt-get install -y less man curl wget git unzip bzip2 vim emacs && \
    apt-get install -y build-essential gcc clang zlib1g-dev libeigen3-dev cmake && \
    apt-get install -y python3 python3-dev python3-pip python3-venv python3-setuptools && \
    touch /root/.bashrc

# install anaconda
RUN wget -P /opt https://repo.anaconda.com/archive/Anaconda3-2020.02-Linux-x86_64.sh && \
    bash /opt/Anaconda3-2020.02-Linux-x86_64.sh -b -p /opt/anaconda3 && \
    rm /opt/Anaconda3-2020.02-Linux-x86_64.sh && \
    echo "export PATH=/opt/anaconda3/bin:$PATH" >> /root/.bashrc && \
    . /root/.bashrc && \
    conda init

# install server env
RUN apt-get update && \
    apt-get install -y mysql-server golang-go

# socketエラーの対策
RUN mkdir -p /var/run/mysqld/ && \
    chmod 777 /var/run/mysqld/

# install angular
RUN apt install -y nodejs npm && \
    npm install n -g && \
    n stable && \
    . /root/.bashrc && \
    echo n | npm install -g --silent @angular/cli

# キャッシュ削除
RUN apt-get clean && \
    apt-get autoclean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN git clone https://github.com/tohsato-lab/process-manager.git && \
    go get github.com/go-sql-driver/mysql && \
    go get github.com/gorilla/websocket

# 世界時間情報
ADD ./data/zoneinfo.tar.gz /

WORKDIR /process-manager

RUN service mysql restart && \
    mysql -u root < ./setup/data/init.sql

RUN cd ./app && \
    npm install

RUN . /root/.bashrc && \
    conda config --set auto_activate_base false && \
    conda env create -f=./setup/data/conda_pytorch.yml && \
    conda env create -f=./setup/data/conda_keras.yml

RUN chmod 777 ./launch.sh

