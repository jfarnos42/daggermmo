# Image names
DEV_IMAGE_NAME=daggerfall-server-dev
RELEASE_IMAGE_NAME=daggerfall-server
MODE ?= dev

# Paths and configuration
ifeq ($(MODE),release)
	IMAGE_NAME=$(RELEASE_IMAGE_NAME)
	CONTAINER_WORKDIR=/home/daggeruser
	CODE_MOUNT=
	DB_PATH=/home/daggeruser/daggerfall.db
else
	IMAGE_NAME=$(DEV_IMAGE_NAME)
	CONTAINER_WORKDIR=/app
	CODE_MOUNT=-v $(PWD):/app
	DB_PATH=$(PWD)/Local/daggerfall.db
endif

DB_VOLUME=daggerfall-db

# Docker run base command
DOCKER_RUN=docker run --rm -it \
	$(CODE_MOUNT) \
	-v $(DB_VOLUME):/home/daggeruser \
	-w $(CONTAINER_WORKDIR) \
	$(IMAGE_NAME)

.PHONY: dev release local test initdb run run-release create-volume clean fclean

# Build images
dev:
	docker build -f Dockerfile -t $(DEV_IMAGE_NAME) .

release:
	docker build -f Dockerfile -t $(RELEASE_IMAGE_NAME) --target runtime .

# Build locally
local:
	mkdir -p Local
	go build -o Local/daggerfall ./cmd/server/main.go

# Run tests
test:
	$(DOCKER_RUN) go test ./... -v

# Create DB volume
create-volume:
	docker volume create $(DB_VOLUME)

# Initialize DB
initdb: create-volume
ifeq ($(MODE),release)
	$(DOCKER_RUN) ./daggerfall -initdb -dbpath=$(DB_PATH)
else
	$(DOCKER_RUN) go run cmd/server/main.go -initdb -dbpath=$(DB_PATH)
endif

# Run in Docker
run: create-volume
	docker run -it \
		$(CODE_MOUNT) \
		-v $(DB_VOLUME):/home/daggeruser \
		-w $(CONTAINER_WORKDIR) \
		-p 8080:8080 \
		-p 7777:7777 \
		$(IMAGE_NAME)

# Shortcut for release mode run
run-release: MODE=release
run-release: run

# Clean local build
clean:
	rm -rf Local/daggerfall Local/daggerfall.db

# Full clean including Docker artifacts
fclean: clean
	docker ps -aq | xargs -r docker stop
	docker ps -aq | xargs -r docker rm
	docker images -aq | xargs -r docker rmi -f
	docker volume ls -q | xargs -r docker volume rm
 