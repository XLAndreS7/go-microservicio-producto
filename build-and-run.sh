#!/bin/bash

# Nombre del servicio y etiqueta
IMAGE_NAME="xlandres/producto-service"
TAG="latest"

echo "🔄 Construyendo imagen de Docker..."
docker build -t ${IMAGE_NAME}:${TAG} .

if [ $? -ne 0 ]; then
    echo "❌ Error al construir la imagen. Abortando."
    exit 1
fi

echo "✅ Imagen construida correctamente: ${IMAGE_NAME}:${TAG}"

echo "🧼 Eliminando contenedores existentes (si los hay)..."
docker-compose down

echo "🚀 Levantando contenedores con Docker Compose..."
docker-compose up -d

if [ $? -eq 0 ]; then
    echo "🎉 Microservicio y MongoDB ejecutándose correctamente."
    echo "🌐 Accede al microservicio en http://localhost:8084"
else
    echo "❌ Error al levantar los contenedores."
    exit 1
fi

