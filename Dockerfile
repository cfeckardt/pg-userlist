FROM golang:1.9-alpine
MAINTAINER Patrik Sundberg <patrik.sundberg@gmail.com>

RUN apk add --no-cache git && \
    adduser -u 100 -S output-user && \
    mkdir -p /output && \
    chown output-user /output

VOLUME /output

WORKDIR /go/src/app
COPY . /go/src/app/

RUN go-wrapper download && go-wrapper install

USER output-user
