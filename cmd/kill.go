package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/NoaTamburrini/portman/internal/process"
	"github.com/NoaTamburrini/portman/internal/scanner"
)

func executeKill() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: portman kill <port>")
		os.Exit(1)
	}

	portNum, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid port number: %s\n", os.Args[2])
		os.Exit(1)
	}

	if portNum < 1 || portNum > 65535 {
		fmt.Fprintf(os.Stderr, "Port number must be between 1 and 65535\n")
		os.Exit(1)
	}

	// Scan to find the port
	ports, err := scanner.ScanPorts()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning ports: %v\n", err)
		os.Exit(1)
	}

	port := scanner.FindByPort(ports, portNum)
	if port == nil {
		fmt.Printf("No process found on port %d\n", portNum)
		os.Exit(1)
	}

	fmt.Printf("Killing process on port %d (PID: %d, Process: %s)...\n",
		port.Number, port.PID, port.ProcessName)

	result := process.KillProcess(port.PID)
	if result.Success {
		fmt.Printf("✓ %s\n", result.Message)
	} else {
		fmt.Fprintf(os.Stderr, "✗ %s\n", result.Message)
		os.Exit(1)
	}
}
