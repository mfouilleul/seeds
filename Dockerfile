FROM golang

COPY . /go/src/github.com/mfouilleul/seeds
WORKDIR /go/src/github.com/mfouilleul/seeds

RUN go build -o /go/bin/seeds
