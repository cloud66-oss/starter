#!/usr/bin/env bash

apt-get update
apt-get install -y docker golang pkg-config docker.io

echo "export GOPATH=/var/go" >> ~/.bashrc
