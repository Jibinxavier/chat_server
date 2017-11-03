#!/bin/bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o chatserver_docker chatserver_docker.go
docker build -t chatserver .
