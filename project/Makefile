FRONT_END_BINARY = frontApp
BROKER_BINARY = brokerApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images ..."
	docker compose up -d
	@echo "Starting Docker images ..."

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_docker
	@echo "Stopping docker images (if running ...)"
	docker compose down
	@echo "Building (when required) and starting Docker images ..."
	docker compose up -d
	@echo "Docker Images built and started!"

# down: stop docker compose 
down: 
	@echo "Stopping docker compose ..."
	docker-compose down
	@echo "Done!"

# build_docker: builds the broker binary as a linux executable
build_docker: 
	@echo "Building broker binary ..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"
