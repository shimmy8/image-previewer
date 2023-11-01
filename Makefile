BIN := "./bin/calendar"
DOCKER_IMG="image-previewer-img"
SERVICE_NAME="image-previewer"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/image-previewer

run: build-img
	docker-compose up -d  $(SERVICE_NAME)

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .


.PHONY: build run build-img test lint
