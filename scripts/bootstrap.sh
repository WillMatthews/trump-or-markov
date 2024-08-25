#!/bin/bash

# add to your .bashrc or .zshrc (I'm not going to be rude and fuck with your shell)
# export GOPATH=$HOME/go
# export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

set -e -o pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

check_su(){
    if [ "$EUID" -ne 0 ]; then
        echo -e "${RED}Please run as root.${NC}"
        exit
    fi
}

install_golang() {
    echo "Installing Golang."
    rm -rf /usr/local/go
    curl -sL https://go.dev/dl/go1.23.0.linux-amd64.tar.gz | tar -C /usr/local -xz
}

#
#
#

main() {
    check_su
    echo "Installing dependencies."
    install_golang
    echo -e "${GREEN}Dependencies installed successfully.${NC}"
}

main
