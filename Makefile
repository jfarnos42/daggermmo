DEV_IMAGE_NAME=daggerfall-server-dev
RELEASE_IMAGE_NAME=daggerfall-server
DB_VOLUME=daggerfall-db

# Paths
LOCAL_DIR=$(PWD)/Local
DB_PATH_DEV=$(LOCAL_DIR)/daggerfall.db
DB_PATH_REL=/home/daggeruser/daggerfall.db

# Build dev image
build-dev:
	docker build -f Dockerfile -t $(DEV_IMAGE_NAME) .

# Build release image
build-release:
	docker build -f Dockerfile -t $(RELEASE_IMAGE_NAME) --target runtime .

# Create Local directory if missing
local-dir:
	mkdir -p $(LOCAL_DIR)

# Create Docker volume for DB
create-volume:
	docker volume create $(DB_VOLUME) || true

# Initialize DB in dev mode
initdb-dev: local-dir create-volume
	docker run --rm -it \
		-v $(PWD):/app \
		-v $(DB_VOLUME):/home/daggeruser \
		-w /app \
		$(DEV_IMAGE_NAME) \
		go run cmd/server/main.go -initdb -dbpath=$(DB_PATH_REL)

# Initialize DB in release mode
initdb-release: create-volume
	docker run --rm -it \
		-v $(DB_VOLUME):/home/daggeruser \
		-w /home/daggeruser \
		$(RELEASE_IMAGE_NAME) \
		./daggerfall -initdb -dbpath=$(DB_PATH_REL)

# Run dev container
run-dev: local-dir create-volume
	docker run -it \
		-v $(PWD):/app \
		-v $(DB_VOLUME):/home/daggeruser \
		-w /app \
		-p 8080:8080 -p 7777:7777 \
		$(DEV_IMAGE_NAME)

# Run release container
run-release: create-volume
	docker run -it \
		-v $(DB_VOLUME):/home/daggeruser \
		-w /home/daggeruser \
		-p 8080:8080 -p 7777:7777 \
		$(RELEASE_IMAGE_NAME)

# Clean local build files
clean:
	rm -rf Local/daggerfall Local/

# Full clean
fclean: clean
	docker ps -aq | xargs -r docker stop
	docker ps -aq | xargs -r docker rm
	docker images -aq | xargs -r docker rmi -f
	docker volume ls -q | xargs -r docker volume rm
