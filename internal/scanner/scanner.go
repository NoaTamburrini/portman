package scanner

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// ScanPorts scans for all active ports on the system
func ScanPorts() ([]Port, error) {
	switch runtime.GOOS {
	case "darwin", "linux":
		return scanPortsUnix()
	case "windows":
		return scanPortsWindows()
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// scanPortsUnix uses lsof to scan ports on macOS and Linux
func scanPortsUnix() ([]Port, error) {
	cmd := exec.Command("lsof", "-i", "-P", "-n")
	output, err := cmd.Output()
	if err != nil {
		// lsof returns non-zero exit code when no processes found
		if exitErr, ok := err.(*exec.ExitError); ok {
			if len(exitErr.Stderr) == 0 {
				return []Port{}, nil
			}
		}
		return nil, fmt.Errorf("failed to execute lsof: %w", err)
	}

	return parseUnixOutput(string(output))
}

// parseUnixOutput parses the output from lsof
func parseUnixOutput(output string) ([]Port, error) {
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return []Port{}, nil
	}

	portMap := make(map[string]Port) // Use map to deduplicate

	for _, line := range lines[1:] { // Skip header
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		processName := fields[0]
		pidStr := fields[1]
		protocol := strings.ToLower(fields[7])
		address := fields[8]

		// Parse PID
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		// Extract port number from address
		var port int
		if strings.Contains(address, ":") {
			parts := strings.Split(address, ":")
			if len(parts) >= 2 {
				portStr := parts[len(parts)-1]
				// Remove any trailing characters like (LISTEN)
				portStr = strings.Split(portStr, "(")[0]
				port, err = strconv.Atoi(portStr)
				if err != nil {
					continue
				}
			}
		}

		if port == 0 {
			continue
		}

		// Get command (rest of the line)
		command := processName
		if len(fields) > 9 {
			command = strings.Join(fields[9:], " ")
		}

		// Create unique key for deduplication
		key := fmt.Sprintf("%s-%d-%d", protocol, port, pid)

		portMap[key] = Port{
			Number:      port,
			PID:         pid,
			ProcessName: processName,
			Command:     command,
			Protocol:    protocol,
		}
	}

	// Convert map to slice
	ports := make([]Port, 0, len(portMap))
	for _, port := range portMap {
		ports = append(ports, port)
	}

	return ports, nil
}

// scanPortsWindows uses netstat to scan ports on Windows
func scanPortsWindows() ([]Port, error) {
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute netstat: %w", err)
	}

	return parseWindowsOutput(string(output))
}

// parseWindowsOutput parses the output from netstat on Windows
func parseWindowsOutput(output string) ([]Port, error) {
	lines := strings.Split(output, "\n")
	if len(lines) < 4 {
		return []Port{}, nil
	}

	portMap := make(map[string]Port)

	for _, line := range lines[4:] { // Skip headers
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		protocol := strings.ToLower(fields[0])
		localAddress := fields[1]
		pidStr := fields[len(fields)-1]

		// Parse PID
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		// Extract port from local address
		parts := strings.Split(localAddress, ":")
		if len(parts) < 2 {
			continue
		}

		port, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil || port == 0 {
			continue
		}

		// Get process name from PID (Windows specific)
		processName := getProcessNameWindows(pid)

		key := fmt.Sprintf("%s-%d-%d", protocol, port, pid)

		portMap[key] = Port{
			Number:      port,
			PID:         pid,
			ProcessName: processName,
			Command:     processName,
			Protocol:    protocol,
		}
	}

	ports := make([]Port, 0, len(portMap))
	for _, port := range portMap {
		ports = append(ports, port)
	}

	return ports, nil
}

// getProcessNameWindows gets the process name from PID on Windows
func getProcessNameWindows(pid int) string {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/FO", "CSV", "/NH")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	fields := strings.Split(strings.TrimSpace(string(output)), ",")
	if len(fields) > 0 {
		return strings.Trim(fields[0], "\"")
	}

	return "unknown"
}

// FindByPort finds a port by its port number
func FindByPort(ports []Port, portNum int) *Port {
	for _, p := range ports {
		if p.Number == portNum {
			return &p
		}
	}
	return nil
}
