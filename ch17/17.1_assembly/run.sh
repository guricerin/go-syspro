#!/usr/bin/env bash
set -ex

gcc -S main.c
gcc main.s
readelf --syms a.out
readelf --file-header a.out
