APP_NAME := ibcmon
DOCKER_REPO := dlvlabs/$(APP_NAME)
VERSION := $(shell git describe --tags --abbrev=0)

PLATFORMS := linux/amd64,linux/arm64

run:
	go run main.go -config config.toml

docker-build:
	docker build -t $(APP_NAME):latest .

docker-push:
	@echo "Building image for all platforms: $(PLATFORMS)"
	docker buildx inspect multiarch-builder >/dev/null 2>&1 || docker buildx create --use --name multiarch-builder --bootstrap
	docker buildx use multiarch-builder
	docker buildx build \
		--platform $(PLATFORMS) \
		--tag ghcr.io/$(DOCKER_REPO):latest \
		--tag ghcr.io/$(DOCKER_REPO):$(VERSION) \
		--file Dockerfile \
		--push .
