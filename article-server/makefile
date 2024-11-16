IMAGE_NAME = go-http-server
CONTAINER_NAME = go-http-server-container

build:
	docker build -t $(IMAGE_NAME) .

run: stop
	docker run -d --name $(CONTAINER_NAME) -p 9999:9999 $(IMAGE_NAME)

stop:
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

logs:
	docker logs $(CONTAINER_NAME)

clean:
	docker rm -f $(docker ps -aq) || true
	docker rmi -f $(docker images -q -f "dangling=true") || true

test:
	curl -i http://localhost:9999

.PHONY: build run stop logs clean test
