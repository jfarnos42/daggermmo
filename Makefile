# Nombre de la imagen
IMAGE_NAME=daggerfall-server-dev

# Directorio local donde se almacenará la base de datos
DB_VOLUME=$(PWD)/dbdata

# Ruta del código dentro del contenedor
CONTAINER_WORKDIR=/app

# Comando base de docker run
DOCKER_RUN=docker run --rm -it \
	-v $(PWD):$(CONTAINER_WORKDIR) \
	-v $(DB_VOLUME):/home/daggeruser \
	-w $(CONTAINER_WORKDIR) \
	$(IMAGE_NAME)

.PHONY: all build test initdb run

# Compila la imagen de desarrollo
build:
	docker build -t $(IMAGE_NAME) .

# Corre todos los tests del proyecto
test:
	$(DOCKER_RUN) go test ./... -v

# Inicializa la base de datos (crea daggerfall.db en dbdata/)
initdb:
	$(DOCKER_RUN) go run cmd/server/main.go -initdb

# Ejecuta el servidor (sin --rm para mantenerlo activo)
run:
	docker run -it \
		-v $(PWD):$(CONTAINER_WORKDIR) \
		-v $(DB_VOLUME):/home/daggeruser \
		-w $(CONTAINER_WORKDIR) \
		-p 8080:8080 \
		-p 7777:7777 \
		$(IMAGE_NAME) go run cmd/server/main.go
