FROM golang:alpine

LABEL maintainer "John Arroyo <john.arroyo@gmail.com>"

RUN go get github.com/arroyo/grutil

# Config
RUN mkdir ~/.grutil
COPY config.yaml ~/.grutil

# Templates
COPY templates /go/templates
