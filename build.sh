#!/bin/bash
go get -u -v
export GOARCH="amd64"
export GOOS="windows"
export CGO_ENABLED=1
go build -v -ldflags="-s -w"
export GOOS="linux"
export CGO_ENABLED=0
go build -v -ldflags="-s -w"
mkdir -p usr/lib64/nagios/plugins etc/icinga2
cp check_f5_throughput usr/lib64/nagios/plugins/
cp sample_config.yaml etc/icinga2/check_f5_throughput.yaml
tar -czvf check_f5_throughput.tar.gz usr etc
rm -r usr etc
