#!/bin/sh
rm -rf ./_bin

mkdir -p ./_bin

# mac
echo "\n\n*** MAC:"
go build -o _bin/gamma-cli
tar -zcvf ./_bin/gammacli-mac-v1.2.3.tar.gz ./_bin/gamma-cli
openssl sha256 ./_bin/gammacli-mac-v1.2.3.tar.gz

# linux
echo "\n\n*** LINUX:"
GOOS=linux go build -o _bin/gamma-cli
tar -zcvf ./_bin/gammacli-linux-v1.2.3.tar.gz ./_bin/gamma-cli
openssl sha256 ./_bin/gammacli-linux-v1.2.3.tar.gz
