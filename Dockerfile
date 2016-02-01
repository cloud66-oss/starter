#use the ruby base image
FROM golang:1.5
MAINTAINER Daniël van Gils

#get alll the stuff
RUN go get github.com/bugsnag/bugsnag-go
RUN go get github.com/mgutz/ansi
RUN go get github.com/hashicorp/go-version
RUN go get github.com/mitchellh/gox

#copy the source files
RUN mkdir -p /usr/local/go/src/github.com/cloud66/starter
ADD . /usr/local/go/src/github.com/cloud66/starter

#switch to our app directory
WORKDIR /usr/local/go/src/github.com/cloud66/starter