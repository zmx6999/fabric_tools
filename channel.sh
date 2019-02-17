#!/bin/bash

if [ $# -lt 1 ]; then echo `usage: chaincode.sh CONFIG_PATH`; exit 1; fi

if [ ! -e channel.go ]; then echo "channel.go not found"; exit 1; fi

go run channel.go $1
if [ $? -ne 0 ]; then exit 1; fi

chmod +x _channel.sh
./_channel.sh
if [ $? -ne 0 ]; then exit 1; fi
