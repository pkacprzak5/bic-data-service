run: build
	@./bin/api

build:
	@go build -o ./bin/api ./cmd

test:
	@go test -v ./...

test-unit:
	@go test -v -short ./...

test-integration:
	@go test -v -tags=integration ./...