# Makefile for the Go API Gateway project
.PHONY: help run build test tidy docker-build
APP_NAME=api-gateway
CMD_PATH=./cmd/api
BINARY_NAME=api-gateway
DOCKER_IMAGE_NAME=template-go-api
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  run           Run the application locally for development"
	@echo "  build         Build the application for Linux AMD64"
	@echo "  test          Run all tests"
	@echo "  tidy          Tidy go.mod and go.sum"
	@echo "  docker-build  Build the Docker image"

run: tidy
	@echo "Running the application..."
	go run $(CMD_PATH)

build: tidy
	@echo "Building binary for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/$(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete: ./bin/$(BINARY_NAME)"

tidy:
	@echo "Tidying module files..."
	go mod tidy

docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE_NAME):latest .
	@echo "Docker image built: $(DOCKER_IMAGE_NAME):latest"
	
clean:
	@echo "Cleaning up..."
	rm -f ${BINARY_NAME}

swaginit:
	@echo "Swagger install dependencies..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go get github.com/swaggo/files
	go get github.com/swaggo/gin-swagger
	go get github.com/swaggo/swag
	go mod tidy

swagdoc:
	@echo "Swagger init doc..."
	swag init -g ./cmd/server/main.go

testinit:
	@echo "Testify install dependencies..."
	go get github.com/stretchr/testify

test:
	@echo "Testify running..."
	go test ./tests/... -v