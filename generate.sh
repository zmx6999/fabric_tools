#!/bin/bash

if [ ! -e generateConfig.go ]; then echo "generateConfig.go not found"; exit 1; fi

go run generateConfig.go
if [ $? -ne 0 ]; then exit 1; fi

chmod +x _generate.sh
./_generate.sh
