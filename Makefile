GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=cyber_bed
LINTER=golangci-lint

DOCKER_IMAGES=$(docker images -aq)
DOCKER_VOLUMES=$(docker volume ls -q)

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

lint:
	$(LINTER) run

docker_up: docker_clean_full
	docker compose up -d

docker_clean_full:
	docker compose down
	# docker rmi -f $(DOCKER_IMAGES)  
	# docker volume rm $(DOCKER_VOLUMES)

docker_logs:
	docker compose logs
