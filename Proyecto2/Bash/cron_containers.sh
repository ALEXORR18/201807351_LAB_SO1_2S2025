#!/bin/bash
# Archivo: cron_containers.sh
# Propósito: Crear contenedores aleatorios para el proyecto SO1

# Definir imágenes
HIGH_CONSUME_IMAGES=("so1_high_cpu" "so1_high_mem")
LOW_CONSUME_IMAGES=("so1_low")

# Cantidad de contenedores a crear
NUM_CONTAINERS=10

# Prefijo de nombre
PREFIX="so1_test"

# Contador para nombres aleatorios
COUNTER=1

# Crear contenedores
for i in $(seq 1 $NUM_CONTAINERS); do
    # Elegir aleatoriamente alto o bajo consumo
    if [ $((RANDOM % 2)) -eq 0 ]; then
        IMAGE=${HIGH_CONSUME_IMAGES[$((RANDOM % ${#HIGH_CONSUME_IMAGES[@]}))]}
    else
        IMAGE=${LOW_CONSUME_IMAGES[$((RANDOM % ${#LOW_CONSUME_IMAGES[@]}))]}
    fi

    # Generar nombre único
    NAME="${PREFIX}_$RANDOM"

    # Crear contenedor en segundo plano
    docker run -d --name "$NAME" "$IMAGE" > /dev/null 2>&1

    COUNTER=$((COUNTER + 1))
done

# Limpiar contenedores parados antiguos (opcional)
docker container prune -f > /dev/null 2>&1
# Listar contenedores activos
echo "Contenedores activos:"
docker ps --format "table {{.Names}}\t{{.ID}}\t{{.Status