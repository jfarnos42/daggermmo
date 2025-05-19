#!/bin/bash
set -e

echo "Parando todos los contenedores..."
docker ps -aq | xargs -r docker stop

echo "Borrando todos los contenedores..."
docker ps -aq | xargs -r docker rm

echo "Borrando todas las imágenes (¡CUIDADO!)..."
docker images -aq | xargs -r docker rmi -f

