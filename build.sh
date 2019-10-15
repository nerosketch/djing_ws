#!/bin/bash

GOPATH=`pwd`:`pwd`/src

#go build -ldflags "-s -w" main.go
go build main.go
