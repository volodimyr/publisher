GOCMD=go
GOTEST=$(GOCMD) test

test:
	${GOTEST} -v -cover ./...
run:
	${GOCMD} run cmd/main.go
docker-build:
	docker image build -t publisher:latest .
docker-run:
	docker container run -d -p 8080:8080 --name publisher publisher
docker-stop:
	docker stop publisher && docker rm publisher
docker-build-run:
	docker image build -t publisher:latest . && docker container run -d -p 8080:8080 --name publisher publisher