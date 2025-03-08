package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

const ProducerURL = "http://localhost:5000/api/producer/v1/stats"

type SystemStats struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage uint64  `json:"memory_usage"`
	DiskUsage   uint64  `json:"disk_usage"`
	Timestamp   string  `json:"timestamp"`
}

func GetSystemStats() (*SystemStats, error) {
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	memStats, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	diskStats, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	stats := &SystemStats{
		CPUUsage:    cpuPercent[0],
		MemoryUsage: memStats.Used,
		DiskUsage:   diskStats.Used,
		Timestamp:   time.Now().Format(time.RFC3339),
	}
	return stats, nil
}

func SendStats(stats *SystemStats) error {
	data, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	resp, err := http.Post(ProducerURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to send stats:%d", resp.StatusCode)
	}
	return nil
}

func StartCronJob() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stats, err := GetSystemStats()
		if err != nil {
			fmt.Println("Error getting system stats:", err)
			continue
		}
		err = SendStats(stats)
		if err != nil {
			log.Println("Error sending stats to producer", err)
		}
	}
}

func main() {
	fmt.Println("Starting system stats cron job...")
	StartCronJob() // Starts the cron job
}
