#!/bin/bash

readonly BIN="a.out"
go build -o ${BIN}
sudo dtruss ./${BIN}