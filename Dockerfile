FROM golang:1.9-alpine
MAINTAINER Patrik Sundberg <patrik.sundberg@gmail.com>

RUN apk add --no-cache git && \
    adduser -u 2323 -S pgpool && \
    mkdir -p /etc/pgpool && \
    chown pgpool /etc/pgpool

VOLUME /etc/pgpool

WORKDIR /go/src/app
COPY . /go/src/app/

RUN go-wrapper download && go-wrapper install

USER pgpool
