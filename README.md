# EventPublisher
Service performs http requests to registered service in case of event published
### Run tests
```sh
$ make test
```
or run
```go test ./...```
### Run server locally
```sh
$ make run
```
It executes ```go run api/publisher.go``` under the hood. Which is the entrypoint for the api.
Control + C if you want to drop server
### Build & Run docker image
```sh
$ make docker-build-run
```
It executes: ```docker image build -t event_publisher:latest . && docker container run -d -p 8080:8080 --name event_publisher event_publisher``` under the hood

or you can use the next commands one by one
```sh
$ make docker-build
```
and
```sh
$ make docker-run
```
If you want to stop and remove image then use next command
```sh
$ make docker-stop
```
It executes: ```docker stop event_publisher && docker rm event_publisher``` under the hood
