#!/usr/bin/env bash

docker run -t -i --name build_temp -v `pwd`:/go/src/github.com/myles-mcdonnell/sendmailserviceproxy golang:1.8.3-alpine3.6 go install github.com/myles-mcdonnell/sendmailserviceproxy/cmd/sendmailserviceproxy-server
docker cp build_temp:/go/bin/sendmailserviceproxy-server ./docker/sendmailserviceproxy-server
docker rm build_temp
docker build -t sendmailserviceproxy ./docker/