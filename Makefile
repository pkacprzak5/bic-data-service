run: build
	@./bin/api

build:
	@go build -o ./bin/api ./cmd

test:
	@go test -v -tags ./...

test-integration: # unit tests + integration tests
	@go test -v -tags=integration ./...

test-end2end: # unit tests + end2end tests
	@go test -v -tags=end2end ./...

# Docker commands
docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

docker-reset-database:
	docker-compose down -v

docker-test: # unit tests + integration tests
	docker-compose -f docker-compose.test.yaml up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.yaml down -v

docker-test-end2end: #unit tests + end2end tests
	docker-compose -f docker-compose.teste2e.yaml up --build --abort-on-container-exit
	docker-compose -f docker-compose.teste2e.yaml down -v