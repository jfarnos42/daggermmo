#!/bin/bash
set -e

echo "Parando todos los contenedores..."
docker ps -aq | xargs -r docker stop

echo "Borrando todos los contenedores..."
docker ps -aq | xargs -r docker rm

echo "Borrando todas las imágenes (¡CUIDADO!)..."
docker images -aq | xargs -r docker rmi -f

echo "Borrando todos los volúmenes (¡CUIDADO!)..."
docker volume ls -q | xargs -r docker volume rm
