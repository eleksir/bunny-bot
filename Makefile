#!/usr/bin/env gmake -f

GOOPTS=CGO_ENABLED=0
BUILDOPTS=-ldflags="-s -w" -a -gcflags=all=-l -trimpath -buildvcs=false

all: clean build

build:
	${GOOPTS} go build ${BUILDOPTS} -o bunny-bot ./cmd/bunny-bot

clean:
	rm -rf bunny-bot

upgrade:
	rm -rf vendor
	go get -u -t -tool ./...
	go mod tidy
	go mod vendor

# vim: set ft=make noet ai ts=4 sw=4 sts=4:
