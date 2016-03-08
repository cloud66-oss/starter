#use the golang base image
FROM golang:1.5
MAINTAINER Daniël van Gils

#get all the go runtime stuff
RUN go get github.com/bugsnag/bugsnag-go
RUN go get github.com/mgutz/ansi
RUN go get github.com/hashicorp/go-version
RUN go get github.com/mitchellh/gox
RUN go get github.com/mitchellh/go-homedir

#gat all the go testing stuff
RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega

#install habitus
RUN curl -L -o /usr/local/bin/habitus https://github.com/cloud66/habitus/releases/download/v0.3/habitus_linux_amd64 
RUN chmod a+x /usr/local/bin/habitus 

#copy the source files
RUN mkdir -p /usr/local/go/src/github.com/cloud66/starter
ADD . /usr/local/go/src/github.com/cloud66/starter

#switch to our app directory
WORKDIR /usr/local/go/src/github.com/cloud66/starter