# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
BINARY_NAME=alertService
BINARY_VERSION=v1.0.1

# docker parameters
DOCKER_IMAGE_NAME=alertService
DOCKER_IMAGE_TAG=$(BINARY_VERSION)
DOCKER_REGISTRY_URL=
DOCKER_USERNAME=
DOCKER_PASSWORD=

# Builds the project for Linux
amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_$(BINARY_VERSION) -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -v ./cmd
	@echo "build amd64 linux success"

arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BINARY_NAME)_$(BINARY_VERSION) -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -v ./cmd
	@echo "build arm64 linux success"

# Builds the project for docker images
docker:
	docker build -t $(DOCKER_REGISTRY_URL)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .
	docker login -u $(DOCKER_USERNAME) -p $(DOCKER_PASSWORD) $(DOCKER_REGISTRY_URL)
	docker push $(DOCKER_REGISTRY_URL)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
	@echo "build docker images success"

.PHONY: all test clean amd64 arm64 docker

all: clean test  clean amd64 arm64 docker

test:
	$(GOTEST) -v ./...
	@echo "test success"

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME_MAC)
	rm -f $(BINARY_NAME)
	docker rmi $(DOCKER_REGISTRY_URL)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)
	@echo "clean success"

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd
	./$(BINARY_NAME)
	@echo "run success"
