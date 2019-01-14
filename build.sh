#!/bin/sh -xe
rm -rf ./_bin

mkdir -p ./_bin
go build -o _bin/gamma-cli

if [ $(uname) == "Darwin" ]; then
    # mac only for now
    tar -zcvf ./_bin/gammacli-mac-v1.2.3.tar.gz ./_bin/gamma-cli
    openssl sha256 ./_bin/gammacli-mac-v1.2.3.tar.gz
fi