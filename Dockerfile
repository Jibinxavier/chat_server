FROM scratch
ADD chatserver_docker /
ENTRYPOINT ["/chatserver_docker"]