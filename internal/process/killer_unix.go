//go:build !windows

package process

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// KillProcess kills a process by PID with graceful fallback
func KillProcess(pid int) KillResult {
	if pid <= 0 {
		return KillResult{
			Success: false,
			Message: "Invalid PID",
		}
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return KillResult{
			Success: false,
			Message: fmt.Sprintf("Process not found: %v", err),
		}
	}

	// Try graceful kill first (SIGTERM)
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		if err.Error() == "os: process already finished" {
			return KillResult{
				Success: true,
				Message: "Process already terminated",
			}
		}

		// Try force kill immediately if SIGTERM fails
		err = process.Signal(syscall.SIGKILL)
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

	// Wait a bit to see if process terminates gracefully
	terminated := waitForTermination(pid, 2*time.Second)

	if !terminated {
		err = process.Signal(syscall.SIGKILL)
		if err != nil {
			return KillResult{
				Success: false,
				Message: fmt.Sprintf("Failed to force kill process: %v", err),
			}
		}

		return KillResult{
			Success: true,
			Message: "Process killed (forced after timeout)",
		}
	}

	return KillResult{
		Success: true,
		Message: "Process terminated gracefully",
	}
}

func waitForTermination(pid int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		process, err := os.FindProcess(pid)
		if err != nil {
			return true
		}

		err = process.Signal(syscall.Signal(0))
		if err != nil {
			return true
		}

		time.Sleep(100 * time.Millisecond)
	}

	return false
}

// IsProcessRunning checks if a process is still running
func IsProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))
	return err == nil
}
