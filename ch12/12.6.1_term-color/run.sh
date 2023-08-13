#!/usr/bin/env bash

set -ex

go run main.go
go run main.go | tee a.out
