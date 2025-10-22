#!/bin/bash

echo "Eliminando contenedores de prueba..."

# Listamos contenedores cuyo nombre contiene "so1_test_"
CONTENEDORES=$(docker ps -a --filter "name=so1_test_" --format "{{.ID}}")

if [ -z "$CONTENEDORES" ]; then
    echo "No hay contenedores de prueba activos."
else
    for c in $CONTENEDORES; do
        echo "Deteniendo contenedor $c..."
        sudo docker stop $c
        echo "Eliminando contenedor $c..."
        sudo docker rm $c
    done
    echo "Contenedores de prueba eliminados."
fi
echo "Lista de contenedores activos después de la eliminación:"
docker ps --format "table {{.Names}}\t{{.ID}}\t{{.Status
