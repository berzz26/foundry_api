APP_NAME=foundry_api

build:
	go build -o bin/$(APP_NAME) .

run:
	go run cmd/api.go

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf bin

deps:
	go mod tidy

.PHONY: build run test fmt vet clean deps