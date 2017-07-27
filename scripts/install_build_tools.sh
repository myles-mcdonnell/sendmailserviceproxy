#!/bin/bash

set -e

currentPath=`pwd`

rm -rf $GOPATH/src/github.com/go-swagger

go get github.com/go-swagger/go-swagger

cd $GOPATH/src/github.com/go-swagger/go-swagger

git checkout tags/0.10.0

go install github.com/go-swagger/go-swagger/cmd/swagger

cd $currentPath

command -v json-refs >/dev/null 2>&1 || {
    sudo npm install -g json-refs@v2
}