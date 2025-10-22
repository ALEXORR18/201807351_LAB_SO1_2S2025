package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

type LoadType int

const (
	LowLoad LoadType = iota
	HighLoad
)

const (
	CPUHighThreshold = 80.0
	MemHighThreshold = 80.0
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

type ActiveContainer struct {
	Name string
	Type LoadType
}

// ---------------------- Contenedores ----------------------

func randomContainerName(prefix string) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%s_%d", prefix, rand.Intn(100000))
}

func createDynamicContainer(loadType LoadType) ActiveContainer {
	prefix := "so1_test"
	var image, cpu, mem string

	if loadType == LowLoad {
		image = "so1_low"
		cpu = "0.5"
		mem = "256m"
	} else {
		if rand.Intn(2) == 0 {
			image = "so1_high_cpu"
		} else {
			image = "so1_high_mem"
		}
		cpu = "1"
		mem = "512m"
	}

	name := randomContainerName(prefix)
	cmd := exec.Command("docker", "run", "-d", "--name", name, "--cpus", cpu, "--memory", mem, image)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Error creando contenedor %s: %v\n%s\n", name, err, string(out))
	} else {
		fmt.Printf("Contenedor %s creado correctamente (%s)\n", name, image)
	}
	return ActiveContainer{Name: name, Type: loadType}
}

func removeContainer(name string) {
	cmd := exec.Command("docker", "rm", "-f", name)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Error eliminando contenedor %s: %v\n%s\n", name, err, string(out))
	} else {
		fmt.Printf("Contenedor %s eliminado correctamente\n", name)
	}
}

func ensureContainerCounts(activeContainers *[]ActiveContainer) {
	lowCount := 0
	highCount := 0
	for _, c := range *activeContainers {
		if strings.Contains(c.Name, "grafana") {
			continue
		}
		if c.Type == LowLoad {
			lowCount++
		} else {
			highCount++
		}
	}
	for lowCount < 3 {
		ac := createDynamicContainer(LowLoad)
		*activeContainers = append(*activeContainers, ac)
		lowCount++
	}
	for highCount < 2 {
		ac := createDynamicContainer(HighLoad)
		*activeContainers = append(*activeContainers, ac)
		highCount++
	}
}

func manageContainersByLoad(currentCPU, currentMem float64, activeContainers *[]ActiveContainer) {
	ensureContainerCounts(activeContainers)
	if currentCPU > CPUHighThreshold || currentMem > MemHighThreshold {
		ac := createDynamicContainer(LowLoad)
		*activeContainers = append(*activeContainers, ac)
	}
	for i := len(*activeContainers) - 1; i >= 0; i-- {
		c := (*activeContainers)[i]
		if c.Type == LowLoad && !strings.Contains(c.Name, "grafana") {
			lowCount := 0
			for _, x := range *activeContainers {
				if x.Type == LowLoad && !strings.Contains(x.Name, "grafana") {
					lowCount++
				}
			}
			if lowCount > 3 {
				removeContainer(c.Name)
				*activeContainers = append((*activeContainers)[:i], (*activeContainers)[i+1:]...)
				break
			}
		}
	}
}

// ---------------------- Monitoreo ----------------------

func obtenerIDsContenedores() map[string]bool {
	out, err := exec.Command("docker", "ps", "--format", "{{.ID}}").Output()
	if err != nil {
		log.Println("Error ejecutando docker ps:", err)
		return nil
	}
	ids := strings.Split(strings.TrimSpace(string(out)), "\n")
	res := make(map[string]bool)
	for _, id := range ids {
		if id != "" {
			res[id] = true
		}
	}
	return res
}

func esProcesoContenedor(pid int32, contIDs map[string]bool) bool {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/cgroup", pid))
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(data), "\n") {
		for id := range contIDs {
			if strings.Contains(line, id) {
				return true
			}
		}
	}
	return false
}

// ---------------------- SQLite ----------------------

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./so1_metrics.db")
	if err != nil {
		log.Fatal(err)
	}

	createSystem := `CREATE TABLE IF NOT EXISTS system_metrics(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp TEXT,
		total_mem_mb REAL,
		free_mem_mb REAL,
		used_mem_mb REAL,
		cpu_perc REAL
	);`
	createContainer := `CREATE TABLE IF NOT EXISTS container_metrics(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp TEXT,
		name TEXT,
		cpu_perc REAL,
		mem_perc REAL,
		status TEXT
	);`

	db.Exec(createSystem)
	db.Exec(createContainer)
	return db
}

func insertSystemMetrics(db *sql.DB, vm *mem.VirtualMemoryStat) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	db.Exec(`INSERT INTO system_metrics(timestamp,total_mem_mb,free_mem_mb,used_mem_mb,cpu_perc)
		VALUES(?,?,?,?,?)`, timestamp, float64(vm.Total)/1024/1024,
		float64(vm.Available)/1024/1024,
		float64(vm.Used)/1024/1024,
		0.0)
}

func insertContainerMetrics(db *sql.DB, containers []ProcessInfo) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	for _, c := range containers {
		db.Exec(`INSERT INTO container_metrics(timestamp,name,cpu_perc,mem_perc,status)
			VALUES(?,?,?,?,?)`, timestamp, c.Name, c.CPUPerc, c.MemPerc, "running")
	}
}

// ---------------------- Main ----------------------

func main() {
	containerTicker := time.NewTicker(20 * time.Second)
	defer containerTicker.Stop()

	activeContainers := []ActiveContainer{}
	db := initDB()
	defer db.Close()

	for {
		select {
		case <-containerTicker.C:
			go func() {
				vm, _ := mem.VirtualMemory()
				currentCPU := 0.0
				currentMem := (float64(vm.Used) / float64(vm.Total)) * 100
				manageContainersByLoad(currentCPU, currentMem, &activeContainers)
			}()
		default:
			vm, err := mem.VirtualMemory()
			if err != nil {
				log.Println("Error obteniendo info de memoria:", err)
				time.Sleep(5 * time.Second)
				continue
			}

			contIDs := obtenerIDsContenedores()
			allProcs := []ProcessInfo{}
			containers := []ProcessInfo{}

			pids, _ := process.Pids()
			for _, pid := range pids {
				proc, err := process.NewProcess(pid)
				if err != nil {
					continue
				}
				name, _ := proc.Name()
				cmdline, _ := proc.Cmdline()
				memInfo, _ := proc.MemoryInfo()
				cpuPerc, _ := proc.CPUPercent()

				p := ProcessInfo{
					PID:     int(pid),
					Name:    name,
					Cmdline: cmdline,
					VSZ:     int(memInfo.VMS / 1024),
					RSS:     int(memInfo.RSS / 1024),
					MemPerc: float64(memInfo.RSS) / float64(vm.Total) * 100,
					CPUPerc: cpuPerc,
				}

				if esProcesoContenedor(pid, contIDs) {
					containers = append(containers, p)
				} else {
					allProcs = append(allProcs, p)
				}
			}

			sysInfo := SystemInfo{
				TotalMemMB: int(vm.Total / 1024 / 1024),
				FreeMemMB:  int(vm.Available / 1024 / 1024),
				UsedMemMB:  int(vm.Used / 1024 / 1024),
				Processes:  allProcs,
				Containers: containers,
			}

			// Archivos JSON
			out, _ := json.MarshalIndent(sysInfo, "", "  ")
			_ = os.WriteFile("sysinfo.json", out, 0644)

			procsOut, _ := json.MarshalIndent(allProcs, "", "  ")
			_ = os.WriteFile("sysinfo_so1_201807351", procsOut, 0644)

			containersOut, _ := json.MarshalIndent(containers, "", "  ")
			_ = os.WriteFile("continfo_so1_201807351", containersOut, 0644)

			// Inserción en SQLite
			insertSystemMetrics(db, vm)
			insertContainerMetrics(db, containers)

			fmt.Println("sysinfo.json y DB actualizados")
			time.Sleep(5 * time.Second)
		}
	}
}
