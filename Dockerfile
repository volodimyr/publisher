FROM golang:latest
WORKDIR /go/src/github.com/volodimyr/event_publisher
ADD . .
WORKDIR /go/src/github.com/volodimyr/event_publisher/cmd
RUN go build -o cmd
ENTRYPOINT ./cmd