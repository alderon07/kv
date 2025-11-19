APP_NAME := kv
CMD_DIR := ./cmd/kv

.PHONY: all build run test fmt vet tidy clean

all: test build

build:
	go build -o bin/$(APP_NAME) $(CMD_DIR)

run:
	go run $(CMD_DIR) $(ARGS)

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -rf bin


