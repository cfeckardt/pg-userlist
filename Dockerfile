FROM golang:1.10-alpine as builder
MAINTAINER Patrik Sundberg <patrik.sundberg@gmail.com>

VOLUME /output

RUN apk add --no-cache git

WORKDIR /go/src/app
COPY . /go/src/app/

RUN go get -v
RUN go install -v

FROM alpine:3.8

RUN adduser -u 2323 -S app && \
    mkdir -p /output && \
    chown app /output

VOLUME /output

USER app

COPY --from=builder /go/bin/app /app

COPY entrypoint.sh /

ENTRYPOINT /entrypoint.sh
