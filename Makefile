
.PHONY: install-tools install-golinter linter linter-docker run stop restart logs

install-tools: install-golinter

install-golinter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2

linter: install-tools ### Run linter
	golangci-lint run ./...

linter-docker:
	docker run -t --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.62.2 golangci-lint run -v

run:
	docker-compose --env-file .env -f .build/docker-compose.yml up -d --build

stop:
	docker-compose --env-file .env -f .build/docker-compose.yml down

restart: 
	docker-compose --env-file .env -f .build/docker-compose.yml restart

logs:
	docker-compose -f .build/docker-compose.yml logs -f