#!/usr/bin/env bash

rm -rf rootfs
rm -f alpine.tar

docker pull alpine
docker run --name alpine alpine
docker export alpine > alpine.tar
docker rm alpine

mkdir rootfs
tar -C rootfs -xvf alpine.tar

go build -o container main.go
