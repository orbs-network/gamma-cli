#!/bin/bash -e

# Cleanup a gamma-server container in case there is one
docker rm -fv orbs-gamma-server orbs-prism
docker pull orbsnetwork/gamma:v1.0.6

GAMMA_CLI_VERSION="v0.7.1"
GAMMA_CLI_URL="https://github.com/orbs-network/gamma-cli/releases/download/$GAMMA_CLI_VERSION/gammacli-linux-x86-64-$GAMMA_CLI_VERSION.tar.gz"

echo "Downloading pre-built gamma-cli ($GAMMA_CLI_VERSION) from it's official GitHub release repository.."
wget $GAMMA_CLI_URL
tar -zxvf "gammacli-linux-x86-64-$GAMMA_CLI_VERSION.tar.gz"
sudo mv _bin/gamma-cli /usr/bin/gamma-cli

echo "gamma-cli and gamma-server successfully installed"

gamma-cli start-local
