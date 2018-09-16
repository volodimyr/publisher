# Publisher
Service performs http requests to registered service in case of event published
Usually service starts at http://localhost:8080.
###### The next endpoints are available:
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
### Run server locally
```sh
$ make run
```
It executes ```go run cmd/main.go``` under the hood. Which is the entry point for the api.
Control + C if you want to drop server
### Build & Run docker image
```sh
$ make docker-build-run
```
It executes: ```docker image build -t publisher:latest . && docker container run -d -p 8080:8080 --name publisher publisher``` under the hood

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
It executes: ```docker stop publisher && docker rm publisher``` under the hood
