#use the golang base image
FROM golang:1.7
MAINTAINER Daniel van Gils

#get all the go crosscompile stuff
RUN go get github.com/mitchellh/gox

#gat all the go testing stuff
RUN go get github.com/tools/godep
RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega

#copy the source files
RUN mkdir -p /usr/local/go/src/github.com/cloud66/starter
ADD . /usr/local/go/src/github.com/cloud66/starter

#testing without git
ADD ./test/node/express_no_git /usr/local/go/src/github.com/cloud66

#switch to our app directory
WORKDIR /usr/local/go/src/github.com/cloud66/starter