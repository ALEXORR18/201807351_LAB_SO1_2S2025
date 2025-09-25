package main

import (
	//"bufio"
	"encoding/json"
	"fmt"
	"os"

	//"path/filepath"
	//"strconv"
	//"strings"
	"time"
)

// Estructura de cada proceso
type ProcessInfo struct {
	PID  int    `json:"PID"`
	Name string `json:"Name"`
	VSZ  uint64 `json:"VSZ"` // en KB
	RSS  uint64 `json:"RSS"` // en KB
}

// Estructura del JSON final
type SysInfo struct {
	Processes []ProcessInfo `json:"Processes"`
}

// Archivo generado por el módulo
const procFile = "/proc/proc_monitor"
const interval = 5 * time.Second

/* func main() {
	fmt.Println("Daemon iniciado: leyendo", procFile)

	for {
		sysInfo := SysInfo{}

		file, err := os.Open(procFile)
		if err != nil {
			fmt.Println("Error al abrir", procFile, ":", err)
			time.Sleep(interval)
			continue
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			// Ejemplo: "PID:1234 Name:containerd-shim CPU:1000"
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}

			// Parsear PID
			pidStr := strings.TrimPrefix(fields[0], "PID:")
			pid, err := strconv.Atoi(pidStr)
			if err != nil {
				continue
			}

			// Parsear nombre
			name := strings.TrimPrefix(fields[1], "Name:")

			// Obtener VSZ y RSS desde /proc/<pid>/statm
			statmPath := filepath.Join("/proc", pidStr, "statm")
			vsz, rss := getMemStat(statmPath)

			proc := ProcessInfo{
				PID:  pid,
				Name: name,
				VSZ:  vsz,
				RSS:  rss,
			}

			sysInfo.Processes = append(sysInfo.Processes, proc)
		}

		file.Close()

		// Guardar JSON en disco
		outFile := "sysinfo.json"
		saveJSON(outFile, sysInfo)

		time.Sleep(interval)
	}
}

// Función para obtener VSZ y RSS desde /proc/<pid>/statm
func getMemStat(path string) (uint64, uint64) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, 0
	}
	parts := strings.Fields(string(data))
	if len(parts) < 2 {
		return 0, 0
	}
	pageSizeKB := uint64(os.Getpagesize() / 1024)
	vsz, _ := strconv.ParseUint(parts[0], 10, 64)
	rss, _ := strconv.ParseUint(parts[1], 10, 64)
	return vsz * pageSizeKB, rss * pageSizeKB
}

// Función para guardar JSON en disco
func saveJSON(path string, data SysInfo) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creando archivo JSON:", err)
		return
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		fmt.Println("Error escribiendo JSON:", err)
	}
} */

func main() {
	fmt.Println("Daemon iniciado: leyendo", procFile)

	for {
		sysInfo := SysInfo{}

		fileData, err := os.ReadFile(procFile)
		if err != nil {
			fmt.Println("Error al leer", procFile, ":", err)
			time.Sleep(interval)
			continue
		}

		// Parsear JSON completo (array)
		var entries []ProcessInfo
		err = json.Unmarshal(fileData, &entries)
		if err != nil {
			fmt.Println("Error parseando JSON del módulo:", err)
			time.Sleep(interval)
			continue
		}

		sysInfo.Processes = entries

		saveJSON("sysinfo.json", sysInfo)
		time.Sleep(interval)
	}
}

func saveJSON(path string, data SysInfo) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creando archivo JSON:", err)
		return
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		fmt.Println("Error escribiendo JSON:", err)
	}
}
