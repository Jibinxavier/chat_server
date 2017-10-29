FROM golang
ADD . ./main

ENTRYPOINT ./main