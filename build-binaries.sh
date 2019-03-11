#!/bin/sh
rm -rf ./_bin

mkdir -p ./_bin

VERSION="v0.6.6"

echo "\n\n*** MAC:"
GOOS=darwin GOARCH=amd64 go build -o _bin/gamma-cli
tar -zcvf ./_bin/gammacli-mac-$VERSION.tar.gz ./_bin/gamma-cli
rm ./_bin/gamma-cli

echo "\n\n*** LINUX (x86-64):"
GOOS=linux GOARCH=amd64 go build -o _bin/gamma-cli
tar -zcvf ./_bin/gammacli-linux-x86-64-$VERSION.tar.gz ./_bin/gamma-cli
rm ./_bin/gamma-cli

echo "\n\n*** LINUX (i386):"
GOOS=linux GOARCH=386 go build -o _bin/gamma-cli
tar -zcvf ./_bin/gammacli-linux-i386-$VERSION.tar.gz ./_bin/gamma-cli
rm ./_bin/gamma-cli

echo "\n\n*** WINDOWS:"
GOOS=windows GOARCH=386 go build -o _bin/gamma-cli.exe
zip -r ./_bin/gammacli-windows-$VERSION.zip ./_bin/gamma-cli.exe
rm ./_bin/gamma-cli.exe

cd ./_bin

openssl sha256 gammacli-mac-$VERSION.tar.gz >> ./checksums.txt
openssl sha256 gammacli-linux-x86-64-$VERSION.tar.gz >> ./checksums.txt
openssl sha256 gammacli-linux-i386-$VERSION.tar.gz >> ./checksums.txt
openssl sha256 gammacli-windows-$VERSION.zip >> ./checksums.txt
