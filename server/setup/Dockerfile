FROM nvidia/cuda:10.1-cudnn7-devel-ubuntu18.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update
RUN apt install -y mysql-server
RUN apt -y install git
RUN apt -y install golang
RUN apt -y install python3
RUN apt -y install python3-pip
RUN apt -y install wget curl
RUN git clone https://github.com/HiroRittsu/process-manager-server.git && cd process-manager-server && go get github.com/go-sql-driver/mysql && go get github.com/gorilla/websocket

CMD ["service mysql start"]
