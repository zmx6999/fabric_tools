#!/bin/bash

function findGoPackageInstalled() {
    r=(`go env | grep GOROOT | sed 's/\"//g' | sed 's/=/\n/'`)
    goroot=${r[1]}
    path1=(`echo "${goroot}/src/${1}" | sed 's/\/\//\//'`)
    r=(`go env | grep GOPATH | sed 's/\"//g' | sed 's/=/\n/'`)
    gopath=${r[1]}
    path2=(`echo "${gopath}/src/${1}" | sed 's/\/\//\//'`)
    if test -d "${path1}" -o -d "${path2}"; then return 0; else return 1; fi
}

if [ $# -lt 1 ]; then echo "usage: docker_compose_cfggen.sh COMPOSE_PROJECT_NAME"; exit 1; fi

if [ ! -e docker_compose_cfggen.go ]; then echo "docker_compose_cfggen.go not found"; exit 1; fi

COMPOSE_PROJECT_NAME=$1

if [ -e .env ]; then rm -f .env; fi
env=$'IMAGE_TAG=latest\nCOMPOSE_PROJECT_NAME='"${COMPOSE_PROJECT_NAME}"
echo "${env}" >> .env

findGoPackageInstalled github.com/sparrc/go-ping
if [ $? -ne 0 ]; then echo "go get github.com/sparrc/go-ping"; go get github.com/sparrc/go-ping; fi
sysctl -w net.ipv4.ping_group_range="0   2147483647"
go run docker_compose_cfggen.go
