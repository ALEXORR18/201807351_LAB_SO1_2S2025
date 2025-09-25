#!/bin/bash

MODULE_NAME="proceso_monitor"
KO_FILE="${MODULE_NAME}.ko"
PROC_FILE="/proc/proc_monitor"

echo "=== Verificando si el módulo $MODULE_NAME está cargado ==="
if lsmod | grep -q "$MODULE_NAME"; then
    echo "Módulo $MODULE_NAME ya cargado. Descargando..."
    sudo rmmod $MODULE_NAME || { echo "Error al descargar $MODULE_NAME"; exit 1; }
    sleep 1
fi

echo "=== Compilando módulo ==="
make || { echo "Error al compilar $KO_FILE"; exit 1; }

echo "=== Cargando módulo $MODULE_NAME ==="
sudo insmod $KO_FILE || { echo "Error al cargar $KO_FILE"; exit 1; }

echo "=== Últimos mensajes del kernel sobre $MODULE_NAME ==="
sudo dmesg | grep "$MODULE_NAME" | tail -n 10

if [ -e "$PROC_FILE" ]; then
    echo "=== Contenido de $PROC_FILE ==="
    cat "$PROC_FILE"
else
    echo "Archivo $PROC_FILE no encontrado."
fi
echo "=== Proceso completado ==="