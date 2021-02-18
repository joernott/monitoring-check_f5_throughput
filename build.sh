#!/bin/bash
BUILDFLAGS='-ldflags="-s -w"'
GOARCH=amd64
GOOS=windows
CGO_ENABLED=1
go build ${BUILDFLAGS}
GOOS=linux
CGO_ENABLED=0
go build ${BUILDFLAGS}
