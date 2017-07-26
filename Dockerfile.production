# use the alpine base image
FROM alpine:3.6
MAINTAINER Cloud 66

RUN apk --update upgrade && apk add curl ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /app
WORKDIR /app

# copy the binary
COPY ./artifacts/compiled/starter /app
COPY ./templates /app/templates
CMD /app/starter -daemon -templates templates -registry true
