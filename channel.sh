#!/bin/bash

if [ ! -e channel.go ]; then echo "channel.go not found"; exit 1; fi

go run channel.go
if [ $? -ne 0 ]; then exit 1; fi

chmod +x _channel.sh
./_channel.sh
if [ $? -ne 0 ]; then exit 1; fi
chmod +x _anchor.sh
./_anchor.sh
if [ $? -ne 0 ]; then exit 1; fi
