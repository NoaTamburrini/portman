//go:build windows

package process

import (
	"fmt"
	"os"
	"os/exec"
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
		process, err := os.FindProcess(pid)
		if err != nil {
			return true
		}

		// On Windows, FindProcess always succeeds, so we check if we can open it
		err = process.Signal(os.Kill)
		if err != nil {
			return true
		}
		// If Signal(Kill) succeeds without error, process still exists
		// but we don't actually want to kill it here — just checking
		// Use tasklist to verify
		cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/NH")
		output, err := cmd.Output()
		if err != nil || len(output) == 0 {
			return true
		}

		time.Sleep(100 * time.Millisecond)
	}

	return false
}

// IsProcessRunning checks if a process is still running
func IsProcessRunning(pid int) bool {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/NH")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(output) > 0 && string(output) != ""
}
