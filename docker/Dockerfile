FROM golang:alpine

LABEL maintainer "John Arroyo <john.arroyo@gmail.com>"

RUN go get github.com/arroyo/grutil

RUN mkdir ~/.grutil

COPY config.yaml ~/.grutil
