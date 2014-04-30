package common

import (
	"os"
	"runtime"
	"time"
)

// startTime represents the application's starting UNIX timestamp
var startTime = time.Now().Unix()

// osInfo represents basic, static information about the host operating system for this process
type osInfo struct {
	Architecture string
	Hostname     string
	NumCPU       int
	PID          int
	Platform     string
}

// Status represents information about the current process, including the basic, static
// information provided by osInfo
type Status struct {
	Architecture string  `json:"architecture"`
	Hostname     string  `json:"hostname"`
	MemoryMB     float64 `json:"memoryMb"`
	NumCPU       int     `json:"numCpu"`
	NumGoroutine int     `json:"numGoroutine"`
	PID          int     `json:"pid"`
	Platform     string  `json:"platform"`
	Uptime       int64   `json:"uptime"`
}

// OSInfo returns information about the host operating system for this process
func OSInfo() (*osInfo, error) {
	// Get system hostname
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	// Return osInfo
	return &osInfo{
		Architecture: runtime.GOARCH,
		Hostname:     hostname,
		NumCPU:       runtime.NumCPU(),
		PID:          os.Getpid(),
		Platform:     runtime.GOOS,
	}, nil
}

// ServerStatus returns information about the current process status
func ServerStatus() (*Status, error) {
	// Retrieve basic OS information
	osStat, err := OSInfo()
	if err != nil {
		return nil, err
	}

	// Get current memory profile
	mem := &runtime.MemStats{}
	runtime.ReadMemStats(mem)

	// Report memory usage in MB
	memMB := float64((float64(mem.Alloc) / 1000) / 1000)

	// Get current uptime
	uptime := time.Now().Unix() - startTime

	// Return status
	return &Status{
		Architecture: osStat.Architecture,
		Hostname:     osStat.Hostname,
		MemoryMB:     memMB,
		NumCPU:       osStat.NumCPU,
		NumGoroutine: runtime.NumGoroutine(),
		PID:          osStat.PID,
		Platform:     osStat.Platform,
		Uptime:       uptime,
	}, nil
}