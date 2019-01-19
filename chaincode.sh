#!/bin/bash

if [ $# -lt 2 ]; then echo `usage: chaincode.sh <option> CONFIG_PATH
options:
    -i instantiate chaincode
    -u upgrade chaincode
`; exit 1; fi

case $1 in
    -i|-u)
    ;;
    *) echo `usage: chaincode.sh <option> CONFIG_PATH
options:
    -i instantiate chaincode
    -u upgrade chaincode
`; exit 1
    ;;
esac

if [ ! -e chaincode.go ]; then echo "chaincode.go not found"; exit 1; fi

go run chaincode.go $1 $2
if [ $? -ne 0 ]; then exit 1; fi

chmod +x _chaincode.sh
./_chaincode.sh
