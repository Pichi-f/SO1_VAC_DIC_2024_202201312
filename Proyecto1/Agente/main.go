package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Estructuras
type Process struct {
	Pid    int     `json:"pid"`
	Name   string  `json:"name"`
	User   int     `json:"user"`
	State  int     `json:"state"`
	Ram    float64 `json:"ram"`
	Father int     `json:"father"`
}

type Cpu struct {
	Usage     float64   `json:"percentage_used"`
	Pc        int       `json:"pc"`
	Processes []Process `json:"tasks"`
}

type Ram struct {
	Total float64 `json:"total_ram"`
	Free  float64 `json:"free_ram"`
	Used  float64 `json:"used_ram"`
	Perc  float64 `json:"percentage_used"`
	Pc    int     `json:"pc"`
}

var (
	cpuData []Cpu
	ramData []Ram
	mutex   sync.Mutex // Para evitar problemas de concurrencia
)

// GET handler para mostrar datos ordenados
func GetDataHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	// Ordenar los datos por porcentaje de uso (mayor a menor)
	sort.Slice(cpuData, func(i, j int) bool {
		return cpuData[i].Usage > cpuData[j].Usage
	})
	sort.Slice(ramData, func(i, j int) bool {
		return ramData[i].Perc > ramData[j].Perc
	})

	// Respuesta JSON
	response := map[string]interface{}{
		"cpu": cpuData,
		"ram": ramData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// postScheduledData mantiene la lógica original
func postScheduledData() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Simulación CPU
			cmdCpu := exec.Command("sh", "-c", "cat /proc/cpu_202201312")
			outCpu, errCpu := cmdCpu.CombinedOutput()
			var cpu_info Cpu = Cpu{Pc: 1}
			if errCpu != nil {
				cpu_info.Usage = 50.0
			} else {
				json.Unmarshal([]byte(outCpu), &cpu_info)
			}

			jsonValue, _ := json.Marshal(cpu_info)
			http.Post("http://35.238.111.219:4000/cpu", "application/json", bytes.NewBuffer(jsonValue))
			// Simulación RAM
			cmdRam := exec.Command("sh", "-c", "cat /proc/ram_202201312")
			outRam, errRam := cmdRam.CombinedOutput()
			var ram_info Ram = Ram{Pc: 1}
			if errRam != nil {
				ram_info = Ram{Total: 16000, Free: 8000, Used: 8000, Perc: 50.0, Pc: 1}
			} else {
				json.Unmarshal([]byte(outRam), &ram_info)
			}

			jsonValue, _ = json.Marshal(ram_info)
			http.Post("http://35.238.111.219:4000/ram", "application/json", bytes.NewBuffer(jsonValue))
		}
	}
}

func main() {
	fmt.Println("Starting server...")
	router := mux.NewRouter().StrictSlash(true)

	// Endpoints
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome!")
	}).Methods("GET")
	router.HandleFunc("/data", GetDataHandler).Methods("GET") // Nuevo endpoint

	// Iniciar go routine
	go postScheduledData()

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	// Iniciar servidor
	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(router)))
}
