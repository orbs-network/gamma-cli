#!/usr/bin/env bash
# This bash script downloads and installs a release pre-built gamma-cli 
# This script assumes you have Docker installed including elevated permissions for your
# user to perform mundane docker tasks such as 'docker ps' (elevated = without the need for sudo)
# Read more at: https://docs.docker.com/v17.12/install/linux/linux-postinstall/
# If you don't have Docker installed, please have a look here on how to install it on Ubuntu Linux:
# https://docs.docker.com/v17.12/install/linux/docker-ce/ubuntu/

docker ps &> /dev/null

DOCKER_TEST_EXITCODE=$?

if [[ $DOCKER_TEST_EXITCODE != 0 ]]; then
    echo "Docker is not properly installed"
    echo "Read more here: https://docs.docker.com/v17.12/install/linux/linux-postinstall/"
    exit 1
fi

GAMMA_CLI_VERSION="v0.7.0"
GAMMA_CLI_URL="https://github.com/orbs-network/gamma-cli/releases/download/$GAMMA_CLI_VERSION/gammacli-linux-x86-64-$GAMMA_CLI_VERSION.tar.gz"

echo "Downloading pre-built gamma-cli ($GAMMA_CLI_VERSION) from it's official GitHub release repository.."
wget $GAMMA_CLI_URL
tar -zxvf gammacli*.tar.gz
sudo mv _bin/gamma-cli /usr/bin/gamma-cli

echo "gamma-cli successfully installed"

exit 0
