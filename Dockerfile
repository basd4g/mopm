FROM golang:1.15-alpine

ENV GO111MODULE=on

RUN apk add --no-cache \
        alpine-sdk \
        bash \
        sudo \
        git

RUN adduser -D -h /home/mopmuser mopmuser \
    && echo "mopmuser:mopmuser" | chpasswd \
    && echo "mopmuser ALL=(ALL)" >> /etc/sudoers

USER mopmuser
WORKDIR /home/mopmuser
COPY ./go.mod /home/mopmuser/go.mod
RUN go mod download
COPY . /home/mopmuser
RUN go build

CMD go test
