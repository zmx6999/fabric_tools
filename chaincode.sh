#!/bin/bash

if [ $# -lt 1 ]; then echo `usage: chaincode.sh <option> [CHAINCODE_NAME]
    -i instantiate chaincode
    -u upgrade chaincode
`; exit 1; fi

if [ ! -e chaincode.go ]; then echo "chaincode.go not found"; exit 1; fi

if [ $# -gt 1 ]; then
    go run chaincode.go $2
    if [ $? -ne 0 ]; then exit 1; fi
else
    go run chaincode.go
    if [ $? -ne 0 ]; then exit 1; fi
fi

chmod +x _chaincode.sh
chmod +x _chaincode_upgrade.sh

case $1 in
    -i) ./_chaincode.sh
    ;;
    -u) ./_chaincode_upgrade.sh
    ;;
    *) echo `usage: chaincode.sh <option> [CHAINCODE_NAME]
    -i instantiate chaincode
    -u upgrade chaincode
`
    ;;
esac
