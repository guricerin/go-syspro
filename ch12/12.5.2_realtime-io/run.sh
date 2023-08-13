#!/usr/bin/env bash
set -ex

go build -o count.out count/count.go
go run main.go
