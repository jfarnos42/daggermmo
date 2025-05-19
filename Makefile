DEV_IMAGE_NAME=daggerfall-server-dev
RELEASE_IMAGE_NAME=daggerfall-server
MODE ?= dev

ifeq ($(MODE),release)
	IMAGE_NAME=$(RELEASE_IMAGE_NAME)
	CONTAINER_WORKDIR=/home/daggeruser
	CODE_MOUNT=
else
	IMAGE_NAME=$(DEV_IMAGE_NAME)
	CONTAINER_WORKDIR=/app
	CODE_MOUNT=-v $(PWD):/app
endif

DB_VOLUME=daggerfall-db

DOCKER_RUN=docker run --rm -it \
	$(CODE_MOUNT) \
	-v $(DB_VOLUME):/home/daggeruser \
	-w $(CONTAINER_WORKDIR) \
	$(IMAGE_NAME)

.PHONY: dev release build test initdb run local initdb-local clean fclean create-volume

dev:
	docker build -f Dockerfile -t $(DEV_IMAGE_NAME) .

release:
	docker build -f Dockerfile -t $(RELEASE_IMAGE_NAME) --target runtime .

create-volume:
	docker volume create $(DB_VOLUME)

test:
	$(DOCKER_RUN) go test ./... -v

initdb: create-volume
ifeq ($(MODE),release)
	$(DOCKER_RUN) ./daggerfall -initdb
else
	$(DOCKER_RUN) go run cmd/server/main.go -initdb
endif

initdb-local: local
	./Local/daggerfall -initdb -dbpath=Local/daggerfall.db

run: create-volume
	docker run -it \
		$(CODE_MOUNT) \
		-v $(DB_VOLUME):/home/daggeruser \
		-w $(CONTAINER_WORKDIR) \
		-p 8080:8080 \
		-p 7777:7777 \
		$(IMAGE_NAME)

local:
	mkdir -p Local
	go build -o Local/daggerfall ./cmd/server/main.go

clean:
	rm -rf Local/daggerfall Local/daggerfall.db

fclean: clean
	docker ps -aq | xargs -r docker stop
	docker ps -aq | xargs -r docker rm
	docker images -aq | xargs -r docker rmi -f
	docker volume ls -q | xargs -r docker volume rm
