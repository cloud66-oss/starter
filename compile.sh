#!/bin/bash 
gox -ldflags "-X main.BUILDDATE=`date -u '+%Y-%m-%d'` -X main.VERSION=1.3.2" -osarch="darwin/386 darwin/amd64 linux/386 linux/amd64 windows/386 windows/amd64" -output="compiled/{{.Dir}}_{{.OS}}_{{.Arch}}"
CGO_ENABLED=0 gox -ldflags "-X main.BUILDDATE=`date -u '+%Y-%m-%d'` -X main.VERSION=1.3.2" -osarch="linux/amd64" -output="compiled/starter"