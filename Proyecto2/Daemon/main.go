package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

type ProcessInfo struct {
	PID     int     `json:"PID"`
	Name    string  `json:"Name"`
	Cmdline string  `json:"Cmdline,omitempty"`
	VSZ     int     `json:"VSZ"`
	RSS     int     `json:"RSS"`
	MemPerc float64 `json:"Memory_Usage"`
	CPUPerc float64 `json:"CPU_Usage"`
}

type SystemInfo struct {
	TotalMemMB int           `json:"TotalMemMB"`
	FreeMemMB  int           `json:"FreeMemMB"`
	UsedMemMB  int           `json:"UsedMemMB"`
	Processes  []ProcessInfo `json:"Processes"`
	Containers []ProcessInfo `json:"Containers"`
}

func main() {
	for {
		content, err := os.ReadFile("/proc/proc_monitor")
		if err != nil {
			log.Println("Error leyendo /proc/proc_monitor:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		var procs []ProcessInfo
		if err := json.Unmarshal(content, &procs); err != nil {
			log.Println("Error parseando JSON del módulo:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		vm, err := mem.VirtualMemory()
		if err != nil {
			log.Println("Error obteniendo info de memoria:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		var allProcs []ProcessInfo
		var containers []ProcessInfo

		for _, p := range procs {
			proc, err := process.NewProcess(int32(p.PID))
			if err != nil {
				continue
			}

			// CMD line
			cmdline, _ := proc.Cmdline()
			p.Cmdline = cmdline

			// Memoria
			memInfo, err := proc.MemoryInfo()
			if err == nil {
				p.RSS = int(memInfo.RSS / 1024)
				p.VSZ = int(memInfo.VMS / 1024)
				p.MemPerc = float64(memInfo.RSS) / float64(vm.Total) * 100
			}

			// CPU
			cpuPerc, err := proc.CPUPercent()
			if err == nil {
				p.CPUPerc = cpuPerc
			}

			allProcs = append(allProcs, p)

			// Filtrar contenedores
			if strings.Contains(strings.ToLower(p.Name), "container") ||
				strings.Contains(strings.ToLower(cmdline), "container") {
				containers = append(containers, p)
			}
		}

		sysInfo := SystemInfo{
			TotalMemMB: int(vm.Total / 1024 / 1024),
			FreeMemMB:  int(vm.Available / 1024 / 1024),
			UsedMemMB:  int(vm.Used / 1024 / 1024),
			Processes:  allProcs,
			Containers: containers,
		}

		out, _ := json.MarshalIndent(sysInfo, "", "  ")
		if err := os.WriteFile("sysinfo.json", out, 0644); err != nil {
			log.Println("Error escribiendo sysinfo.json:", err)
		} else {
			fmt.Println("sysinfo.json actualizado")
		}

		time.Sleep(5 * time.Second)
	}
}
