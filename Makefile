GOCMD=go
GOTEST=$(GOCMD) test

test:
	${GOTEST} -v ./...
run:
	${GOCMD} run api/publisher.go
docker-build:
	docker image build -t event_publisher:latest .
docker-run:
	docker container run -d -p 8080:8080 --name event_publisher event_publisher
docker-stop:
	docker stop event_publisher && docker rm event_publisher
docker-build-run:
	docker image build -t event_publisher:latest . && docker container run -d -p 8080:8080 --name event_publisher event_publisher