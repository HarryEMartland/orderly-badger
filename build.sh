#!/usr/bin/env bash
pwd
ls
echo $GOPATH
docker run -v $GOPATH:/go golang:alpine /bin/sh -c 'cd /go/src/github.com/HarryEMartland/orderly-badger && go get -v -d && go build'
docker build -t harrymartland/orderly-badger .