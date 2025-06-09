APP_NAME := ibcmon
DOCKER_REPO := dlvlabs/$(APP_NAME)
VERSION := $(shell git describe --tags)

PLATFORMS := linux/amd64,linux/arm64

run:
	go run main.go -config config.toml

docker-build:
	docker build -t ibcmon:latest .

docker-push:
	@echo "Building image for all platforms: $(PLATFORMS)"
	docker buildx inspect multiarch-builder >/dev/null 2>&1 || docker buildx create --use --name multiarch-builder --bootstrap
	docker buildx use multiarch-builder
	docker buildx build \
		--platform $(PLATFORMS) \
		--build-arg GOOS=linux \
		--tag $(DOCKER_REPO):$(VERSION) \
		--file Dockerfile \
		--push .
