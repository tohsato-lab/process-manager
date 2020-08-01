#!/bin/bash

git pull origin master

service mysql restart

cd app
bash launch.sh &
cd ../server/src
go run main.go

