#!/bin/bash

git pull origin master

service mysql restart

sleep 5

cd app
bash launch.sh &
cd ../server/src
go run main.go

