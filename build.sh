#!/bin/sh -xe
rm -rf ./_bin

mkdir -p ./_bin
go build -o _bin/gamma-cli