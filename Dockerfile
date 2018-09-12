FROM golang:latest
WORKDIR /go/src/github.com/volodimyr/event_publisher
ADD . .
WORKDIR /go/src/github.com/volodimyr/event_publisher/api
RUN go build -o api
ENTRYPOINT ./api