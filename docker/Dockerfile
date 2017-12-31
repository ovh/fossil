FROM golang:alpine3.7 AS build-env

LABEL MAINTAINER Rachid Zarouali <xinity77@gmail.com>

# install dependencies
RUN apk add --no-cache glide make git

# clone fossil repository and build binary

RUN git clone https://github.com/ovh/fossil.git $GOPATH/src/github.com/ovh/fossil \
    && cd $GOPATH/src/github.com/ovh/fossil \
    && glide install \
    && make release

# final stage
FROM alpine:3.7
COPY --from=build-env /go/src/github.com/ovh/fossil/build/fossil /
EXPOSE 2003
ENTRYPOINT ["/fossil"]