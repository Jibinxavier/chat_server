#!/bin/bash
docker run --rm -it -net=host -e port=$1 chatserver