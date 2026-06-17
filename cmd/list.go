package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/NoaTamburrini/portman/internal/scanner"
	"github.com/charmbracelet/lipgloss"
)

type listOptions struct {
	jsonOutput bool
	all        bool
	port       int // 0 = no filter
}

func executeList() {
	opts := listOptions{}
	for _, arg := range os.Args[2:] {
		switch arg {
		case "--json":
			opts.jsonOutput = true
		case "--all", "-a":
			opts.all = true
		default:
			if n, err := strconv.Atoi(arg); err == nil {
				opts.port = n
			}
		}
	}

	ports, err := scanner.ScanPorts()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning ports: %v\n", err)
		os.Exit(1)
	}

	if opts.port != 0 {
		ports = scanner.FindAllByPort(ports, opts.port)
	}
	if !opts.all {
		ports = scanner.FilterListening(ports)
	}

	sort.Slice(ports, func(i, j int) bool {
		if ports[i].Number != ports[j].Number {
			return ports[i].Number < ports[j].Number
		}
		return ports[i].PID < ports[j].PID
	})

	if opts.jsonOutput {
		if ports == nil {
			ports = []scanner.Port{}
		}
		out, _ := json.Marshal(ports)
		fmt.Println(string(out))
		return
	}

	printListTable(ports)
}

func printListTable(ports []scanner.Port) {
	if len(ports) == 0 {
		fmt.Println("No listening ports found")
		return
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("212")).
		Bold(true)

	header := fmt.Sprintf("%-8s %-9s %-8s %-12s %-20s", "PORT", "PROTOCOL", "PID", "STATE", "PROCESS")
	fmt.Println(headerStyle.Render(header))

	for _, p := range ports {
		state := p.State
		if state == "" {
			state = "-"
		}
		fmt.Printf("%-8d %-9s %-8d %-12s %-20s\n",
			p.Number, p.Protocol, p.PID, state, p.ProcessName)
	}
}
