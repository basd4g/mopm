FROM golang:1.15-alpine
WORKDIR /go-app

ENV GO111MODULE=on

RUN apk add --no-cache \
        alpine-sdk \
        bash \
        git
#RUN echo 'root:root' |chpasswd
#RUN adduser -h /home/mopm mopm \
#    && echo "mopm ALL=(ALL)" >> /etc/sudoers
#    && echo 'mopm:mopm' | chpasswd
#    && echo "mopm ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
#RUN git clone https://github.com/basd4g/mopm.git

COPY . /go-app
RUN go mod download
RUN go build

CMD go test
