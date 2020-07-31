#!/bin/bash

cd app
bash launch.sh &
cd ../server/src
go run main.go

