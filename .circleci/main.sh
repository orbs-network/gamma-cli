#!/bin/bash -xe

PROJ_PATH=`pwd`
GO_VERSION="1.12.6"

# First let's install Go 1.11
echo "Installing Go $GO_VERSION..."
cd /tmp

curl -O https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz
sudo tar -xf go${GO_VERSION}.linux-amd64.tar.gz

# Uninstall older version of Go
sudo rm -rf /usr/local/go
sudo mv go /usr/local

export GOROOT=/usr/local/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

cd $PROJ_PATH

./test.sh
