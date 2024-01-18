GREEN=\033[0;32m
NC=\033[0m

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
BINARY_NAME=alertService
BINARY_VERSION=v1.0.1_release

# docker parameters
DOCKER_IMAGE_NAME=alertService
DOCKER_IMAGE_TAG=$(BINARY_VERSION)
DOCKER_REGISTRY_URL=
DOCKER_USERNAME=
DOCKER_PASSWORD=

# Builds the project for Linux
# Build default is linux amd64
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_$(BINARY_VERSION) -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -v ./cmd
	@printf "$(GREEN)Arch: amd64$(NC)"
	@printf "$(GREEN)OS: linux$(NC)"
	@printf "$(GREEN)build success$(NC)"

build_arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BINARY_NAME)_$(BINARY_VERSION) -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -v ./cmd
	@printf "$(GREEN)Arch: arm64$(NC)"
	@printf "$(GREEN)OS: linux$(NC)"
	@printf "$(GREEN)build success$(NC)"

# Builds the project for docker images
build_docker:
	docker build -t $(DOCKER_REGISTRY_URL)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .
	docker login -u $(DOCKER_USERNAME) -p $(DOCKER_PASSWORD) $(DOCKER_REGISTRY_URL)
	docker push $(DOCKER_REGISTRY_URL)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
	@printf "build docker images success"

.PHONY: all test clean build build_arm64 build_docker

all: clean test  clean build build_arm64 build_docker

test:
	$(GOTEST) -v ./...
	@printf "test success"

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME_MAC)
	rm -f $(BINARY_NAME)
	docker rmi $(DOCKER_REGISTRY_URL)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
	@printf "clean success"

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd
	./$(BINARY_NAME)
	@printf "run success"
