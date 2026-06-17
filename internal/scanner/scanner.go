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

		// Extract the local port. For an established connection the address is
		// "local->remote" (e.g. 192.168.1.45:56060->192.168.1.35:57819) — take the
		// local side so we don't grab the remote port.
		local := address
		if idx := strings.Index(local, "->"); idx != -1 {
			local = local[:idx]
		}

		var port int
		if strings.Contains(local, ":") {
			parts := strings.Split(local, ":")
			portStr := parts[len(parts)-1]
			port, err = strconv.Atoi(portStr)
			if err != nil {
				continue
			}
		}

		if port == 0 {
			continue
		}

		// Connection state, if present, is the next field e.g. "(LISTEN)".
		state := ""
		if len(fields) > 9 {
			state = strings.Trim(fields[9], "()")
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
			State:       state,
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

		// TCP rows carry a State column before the PID (e.g. LISTENING,
		// ESTABLISHED); UDP rows have none. Normalize LISTENING -> LISTEN.
		state := ""
		if strings.HasPrefix(protocol, "tcp") && len(fields) >= 4 {
			state = strings.ToUpper(fields[len(fields)-2])
			if state == "LISTENING" {
				state = "LISTEN"
			}
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
			State:       state,
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

// FindAllByPort finds all processes using a specific port
func FindAllByPort(ports []Port, portNum int) []Port {
	var matches []Port
	for _, p := range ports {
		if p.Number == portNum {
			matches = append(matches, p)
		}
	}
	return matches
}

// FindKillTargets returns the processes to kill for a port: the listening
// (server) sockets if any exist, otherwise all matches. This avoids killing a
// client connection (e.g. a browser tab) that merely shares the port number.
func FindKillTargets(ports []Port, portNum int) []Port {
	matches := FindAllByPort(ports, portNum)

	var listening []Port
	for _, p := range matches {
		if p.IsListening() {
			listening = append(listening, p)
		}
	}

	if len(listening) > 0 {
		return listening
	}
	return matches
}

// FilterListening returns only the listening (server) sockets.
func FilterListening(ports []Port) []Port {
	var listening []Port
	for _, p := range ports {
		if p.IsListening() {
			listening = append(listening, p)
		}
	}
	return listening
}
