FROM nvidia/cuda:10.1-cudnn7-devel-ubuntu18.04

ENV DEBIAN_FRONTEND=noninteractive

RUN add-apt-repository ppa:longsleep/golang-backports && \
    apt update && \
    apt install -y less man curl wget git unzip bzip2 vim emacs && \
    apt install -y build-essential gcc clang zlib1g-dev libeigen3-dev cmake && \
    apt install -y python3 python3-dev python3-pip python3-venv python3-setuptools && \
    # install anaconda
    wget -P /opt https://repo.anaconda.com/archive/Anaconda3-2020.02-Linux-x86_64.sh && \
    bash /opt/Anaconda3-2020.02-Linux-x86_64.sh -b -p /opt/anaconda3 && \
    rm /opt/Anaconda3-2020.02-Linux-x86_64.sh && \
    echo "export PATH=/opt/anaconda3/bin:$PATH" >> ~/.bashrc && \
    . ~/.bashrc && \
    conda init && \
    # install server env
    apt install -y mysql-server golang-go && \
    # tidying up
    apt clean && \
    apt autoclean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN git clone https://github.com/tohsato-lab/process-manager.git && \
    go get github.com/go-sql-driver/mysql && \
    go get github.com/gorilla/websocket

CMD ["service mysql start"]
