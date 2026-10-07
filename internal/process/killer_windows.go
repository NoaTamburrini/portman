//go:build windows

package process

import (
	"encoding/csv"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// KillProcess kills a process by PID using taskkill on Windows
func KillProcess(pid int) KillResult {
	if pid <= 0 {
		return KillResult{
			Success: false,
			Message: "Invalid PID",
		}
	}

	// Try graceful kill first
	cmd := exec.Command("taskkill", "/PID", fmt.Sprintf("%d", pid))
	err := cmd.Run()
	if err == nil {
		terminated := waitForTermination(pid, 2*time.Second)
		if terminated {
			return KillResult{
				Success: true,
				Message: "Process terminated gracefully",
			}
		}
	}

	// Force kill
	cmd = exec.Command("taskkill", "/F", "/PID", fmt.Sprintf("%d", pid))
	err = cmd.Run()
	if err != nil {
		return KillResult{
			Success: false,
			Message: fmt.Sprintf("Failed to kill process: %v", err),
		}
	}

	return KillResult{
		Success: true,
		Message: "Process killed (forced)",
	}
}

func waitForTermination(pid int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if !IsProcessRunning(pid) {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}

	return false
}

// IsProcessRunning checks if a process is still running.
// tasklist prints a localized "no tasks" message instead of empty output when
// nothing matches, so the PID column of the CSV output is checked explicitly.
func IsProcessRunning(pid int) bool {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/FO", "CSV", "/NH")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	reader := csv.NewReader(strings.NewReader(string(output)))
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return false
	}

	want := strconv.Itoa(pid)
	for _, record := range records {
		if len(record) > 1 && strings.TrimSpace(record[1]) == want {
			return true
		}
	}
	return false
}
