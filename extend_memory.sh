#!/bin/bash

set -ev

mkdir /opt/images
dd if=/dev/zero of=/opt/images/swap bs=2048 count=2097152
mkswap /opt/images/swap
swapon /opt/images/swap
