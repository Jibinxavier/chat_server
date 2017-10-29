#!/bin/bash
docker build -t chatserver .
docker run -e port=$1 chatserver