#!/usr/bin/env bash
docker run -v ~/go:/go golang:alpine /bin/sh -c 'cd /go/src/github.com/harryemartland/orderly-badger; go get -v -d; go build'
docker build -t orderly-badger .