FROM golang:1.4-onbuild
MAINTAINER Jeff Anderson <jeff@docker.com>

ENTRYPOINT ["go-wrapper", "run"]
