GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=./cmd/main.go
LINTER=golangci-lint

DOCKER_IMAGES=$(docker images -aq)
DOCKER_VOLUMES=$(docker volume ls -q)

lint:
	$(LINTER) run -c configs/.golangci.yaml

run_api:
	go run cmd/api/main.go -ConfigPath ./configs/app/api/local.yaml

run_notifications:
	go run ./cmd/notifications/notifications.go -ConfigPath ./configs/app/notifications/local.yaml

docker_up: docker_clean_full
	docker compose up -d

docker_logs:
	docker compose logs
