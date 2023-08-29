#!/usr/bin/env bash

set -ex

go build -o check.out check/check.go
go run main.go
