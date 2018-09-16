# EventPublisher
Service performs http requests to registered service in case of event published
Usually service starts at http://localhost:8080.
###### The next endpoint are available:
1. Listener registration
	`POST /listener Body: {"event": "event_name1", "name": "listener_name_1", "address":
	"http://listener.address/handle"}`
2. Listener unregister
	`DELETE /listener/listener_name_1`
3. Publish event
	`POST /publish/{event} Body: json`

### Run tests
```sh
$ make test
```
or run
```go test ./...```
### Run benchmark tests (be careful!)
```sh
$ make btest
```
or run
```	cd api && go test -v -bench=.```
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
