package process

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// KillResult represents the result of a kill operation
type KillResult struct {
	Success bool
	Message string
}

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
		// If SIGTERM fails, might be permission issue or process already dead
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
		// Process didn't terminate, force kill
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

// waitForTermination waits for a process to terminate
func waitForTermination(pid int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		process, err := os.FindProcess(pid)
		if err != nil {
			return true // Process not found means it terminated
		}

		// Signal 0 checks if process exists without actually sending a signal
		err = process.Signal(syscall.Signal(0))
		if err != nil {
			return true // Process doesn't exist anymore
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
