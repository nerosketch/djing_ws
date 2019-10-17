#!/bin/bash
`which protoc` -I=. --go_out=. glob_types/types.proto
