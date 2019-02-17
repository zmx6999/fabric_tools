#!/bin/bash

if [ ! -e generate.go ]; then echo "generate.go not found"; exit 1; fi

go run generate.go
if [ $? -ne 0 ]; then exit 1; fi

chmod +x _generate.sh
./_generate.sh
