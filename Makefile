.PHONY: install-tools install-golinter linter linter-docker test test-coverage generate run stop restart logs

DOCKER_COMPOSE_CMD = docker-compose --env-file .env -f .build/docker-compose.yml

install-tools: install-golinter

install-golinter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2

linter: install-tools
	golangci-lint run ./...

linter-docker:
	docker run -t --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.62.2 golangci-lint run -v

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=cover.out -o=cover.html

generate:
	go generate ./...

run:
	$(DOCKER_COMPOSE_CMD) up -d --build

stop:
	$(DOCKER_COMPOSE_CMD) down

restart: 
	$(DOCKER_COMPOSE_CMD) restart

logs:
	$(DOCKER_COMPOSE_CMD) logs -f