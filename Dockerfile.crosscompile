#use the golang base image
FROM golang:1.8
MAINTAINER Cloud 66

#get all the go crosscompile stuff
RUN go get github.com/mitchellh/gox

#copy the source files
RUN mkdir -p /usr/local/go/src/github.com/cloud66-oss/starter
ADD . /usr/local/go/src/github.com/cloud66-oss/starter

#switch to our app directory
WORKDIR /usr/local/go/src/github.com/cloud66-oss/starter

RUN ./compile.sh
