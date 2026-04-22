APP_NAME := ecommerce-api
CMD := ./cmd/api
BIN  := ./bin/$(APP_NAME)

.PHONY: help tidy run build test test-cover lint fmt vet docker-up docker-down clean

help:
	@echo "Targets: tidy, run, build, test, test-cover, lint, fmt, vet, docker-up, docker-down, clean"

tidy:
	go mod tidy

run:
	go run $(CMD)

build:
	go build -o $(BIN) $(CMD)

test:
	go test ./... -race -count=1

test-cover:
	go test ./... -race -coverprofile=coverage.out
	go tool cover -func=coverage.out

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "install golangci-lint"; exit 1; }
	golangci-lint run ./...

docker-up:
	docker compose up -d

docker-down:
	docker compose down

clean:
	rm -rf bin coverage.out
