# first stage
FROM golang:1.11 as builder
WORKDIR /go/src/github.com/volodimyr/publisher/
COPY . .
WORKDIR /go/src/github.com/volodimyr/publisher/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /publisher/
COPY --from=builder /go/src/github.com/volodimyr/publisher/cmd .
CMD ["./main"]