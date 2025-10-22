#!/bin/bash

MODULE_NAME="proceso_monitor"
KO_FILE="${MODULE_NAME}.ko"
SYS_PROC="/proc/sysinfo_so1_201807351"
CONT_PROC="/proc/continfo_so1_201807351"

echo "=== Verificando si el módulo $MODULE_NAME está cargado ==="
if lsmod | grep -q "$MODULE_NAME"; then
    echo "Módulo $MODULE_NAME ya cargado. Descargando..."
    sudo rmmod $MODULE_NAME || { echo "Error al descargar $MODULE_NAME"; exit 1; }
    sleep 1
fi

echo "=== Compilando módulo ==="
make -C /lib/modules/$(uname -r)/build M=$(pwd) || { echo "Error al compilar $KO_FILE"; exit 1; }

echo "=== Cargando módulo $MODULE_NAME ==="
sudo insmod $KO_FILE || { echo "Error al cargar $KO_FILE"; exit 1; }
sleep 1
echo "Módulo $MODULE_NAME cargado exitosamente."

echo "=== Últimos mensajes del kernel sobre $MODULE_NAME ==="
sudo dmesg | grep "$MODULE_NAME" | tail -n 10

# Mostrar contenido de los archivos /proc
for PROC_FILE in "$SYS_PROC" "$CONT_PROC"; do
    if [ -e "$PROC_FILE" ]; then
        echo "=== Contenido de $PROC_FILE ==="
        sudo cat "$PROC_FILE"
    else
        echo "Archivo $PROC_FILE no encontrado."
    fi
done

echo "=== Proceso completado ==="
# Fin del script