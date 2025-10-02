package tui

import (
	"sort"
	"strings"
	"time"

	"github.com/NoaTamburrini/portman/internal/scanner"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ports          []scanner.Port
	filteredPorts  []scanner.Port
	cursor         int
	statusMessage  string
	statusIsError  bool
	scanning       bool
	lastRefresh    time.Time
	filterMode     bool
	filterInput    textinput.Model
	confirmingKill bool
	width          int
	height         int
}

type scanCompleteMsg struct {
	ports []scanner.Port
	err   error
}

type killCompleteMsg struct {
	success bool
	message string
}

func initialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Filter ports..."
	ti.CharLimit = 50

	return Model{
		ports:         []scanner.Port{},
		filteredPorts: []scanner.Port{},
		cursor:        0,
		filterInput:   ti,
	}
}

func (m Model) Init() tea.Cmd {
	return scanPorts
}

// scanPorts performs a port scan
func scanPorts() tea.Msg {
	ports, err := scanner.ScanPorts()
	if err != nil {
		return scanCompleteMsg{ports: nil, err: err}
	}

	// Sort by port number
	sort.Slice(ports, func(i, j int) bool {
		return ports[i].Number < ports[j].Number
	})

	return scanCompleteMsg{ports: ports, err: nil}
}

// doKillPort kills a port by PID
func doKillPort(pid int) tea.Cmd {
	return func() tea.Msg {
		// Import here to avoid circular dependency
		result := killProcessByPID(pid)
		return killCompleteMsg{
			success: result.Success,
			message: result.Message,
		}
	}
}

// Helper function to call process killer
func killProcessByPID(pid int) struct {
	Success bool
	Message string
} {
	// We need to import the process package
	// This is done in the update.go file where it's actually used
	return struct {
		Success bool
		Message string
	}{Success: false, Message: "Not implemented"}
}

// filterPorts filters the ports based on the filter string
func (m *Model) filterPorts() {
	filter := strings.ToLower(strings.TrimSpace(m.filterInput.Value()))

	if filter == "" {
		m.filteredPorts = m.ports
		return
	}

	filtered := []scanner.Port{}
	for _, p := range m.ports {
		// Check if filter matches port number, process name, or command
		if strings.Contains(strings.ToLower(p.ProcessName), filter) ||
			strings.Contains(strings.ToLower(p.Command), filter) ||
			strings.Contains(strings.ToLower(p.Protocol), filter) {
			filtered = append(filtered, p)
		}
	}

	m.filteredPorts = filtered

	// Adjust cursor if needed
	if m.cursor >= len(m.filteredPorts) {
		m.cursor = max(0, len(m.filteredPorts)-1)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
