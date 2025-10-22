#!/bin/bash

# Número de contenedores de prueba a crear
NUM_CONTENEDORES=3

# Imagen base (ligera) para pruebas
IMAGEN="busybox"

echo "Creando $NUM_CONTENEDORES contenedores de prueba..."

for i in $(seq 1 $NUM_CONTENEDORES); do
    NOMBRE="so1_test_$i"
    # Ejecutamos un contenedor en segundo plano con un sleep largo para que no se cierre
    docker run -d --name "$NOMBRE" "$IMAGEN" sleep 3600
    echo "Contenedor $NOMBRE creado"
done

echo "Contenedores creados. Lista de contenedores activos:"
docker ps --format "table {{.Names}}\t{{.ID}}\t{{.Status}}"
