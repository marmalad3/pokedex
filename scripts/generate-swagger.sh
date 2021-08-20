#!/bin/bash
# This script will generate swagger models & a HTTP client from the
# ./internal/poke/swagger.yaml file

set -ex

docker pull quay.io/goswagger/swagger

swagger="docker run --rm -it  --user $(id -u):$(id -g) -e GOPATH=$HOME/go:/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger"
$swagger generate client -f ./internal/app/pokedex/swagger.yaml -t ./internal/app/pokedex