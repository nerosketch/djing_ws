#!/bin/bash

export GOPATH=`pwd`:`pwd`/src

git submodule init
git submodule update

go get github.com/gorilla/websocket

source ./build.sh
