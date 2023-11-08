GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=./cmd/main.go
LINTER=golangci-lint

DOCKER_IMAGES=$(docker images -aq)
DOCKER_VOLUMES=$(docker volume ls -q)

build:
	$(GOBUILD) -o bin/ $(BINARY_NAME)

lint:
	$(LINTER) run -c configs/.golangci.yaml

docker_up: docker_clean_full
	docker compose up -d

docker_logs:
	docker compose logs
