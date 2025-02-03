run: build
	@./bin/api

build:
	@go build -o ./bin/api ./cmd

test:
	@go test -v ./...

test-integration:
	@go test -v -tags=integration ./...

# Docker commands
docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

docker-reset-database:
	docker-compose down -v

docker-test: # includes all tests
	docker-compose -f docker-compose.test.yaml up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.yaml down -v