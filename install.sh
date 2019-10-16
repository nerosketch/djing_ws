#!/bin/bash

export GOPATH=`pwd`:`pwd`/src

git submodule init
git submodule update

go get github.com/gorilla/websocket
go get -u github.com/golang/protobuf/protoc-gen-go

source ./build.sh
