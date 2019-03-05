#!/bin/bash

set -ev

systemctl stop firewalld.service

yum install docker -y
service docker start
docker version

curl -L https://github.com/docker/compose/releases/download/1.19.0/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
cp /usr/local/bin/docker-compose /usr/bin
docker-compose version

wget https://studygolang.com/dl/golang/go1.11.linux-amd64.tar.gz
tar -C /usr/local -zxvf go1.11.linux-amd64.tar.gz
ln -s /usr/local/go/bin/go /usr/bin/go
go env

yum install git -y
yum install gcc -y

git clone https://github.com/hyperledger/fabric

mkdir -p /root/go/src

go get github.com/hyperledger/fabric/core/chaincode/shim
go get github.com/hyperledger/fabric/protos/peer

cd fabric/scripts
./bootstrap.sh

if [ $# -gt 0 ]; then mkdir -p fabric-samples/${1}/network; fi
