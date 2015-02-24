# https://registry.hub.docker.com/_/golang/
# debian image, GOPATH configured at /go.
FROM golang

ADD . /go/src/github.com/robtuley/httprouter
RUN go get github.com/robtuley/report
RUN go get github.com/robtuley/etcdwatch
RUN go get github.com/robtuley/httpserver

RUN go install github.com/robtuley/httprouter

ENTRYPOINT ["/go/bin/httprouter"]

EXPOSE 8080