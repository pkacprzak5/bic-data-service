run: build
	@./bin/api

build:
	@go build -o ./bin/api ./cmd

test-unit: # unit tests
	@go test -v -tags=unit ./...

test-integration: # integration tests
	@go test -v -tags=integration ./...

test-end2end: # end2end tests
	@go test -v -tags=end2end ./...

# Docker commands
docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

docker-reset-database:
	docker-compose down -v

docker-test-unit: # unit tests
	docker-compose -f docker-compose.test.unit.yaml up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.unit.yaml down -v

docker-test-integration: # integration tests
	docker-compose -f docker-compose.test.integration.yaml up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.integration.yaml down -v

docker-test-end2end: #end2end tests
	docker-compose -f docker-compose.test.e2e.yaml up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.e2e.yaml down -v