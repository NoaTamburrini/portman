package tui

import (
	"fmt"
	"sort"
	"strings"

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
	ti.Width = 40
	ti.PlaceholderStyle = placeholderStyle
	ti.PromptStyle = filterStyle

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
		portNum := fmt.Sprintf("%d", p.Number)
		if strings.Contains(portNum, filter) ||
			strings.Contains(strings.ToLower(p.ProcessName), filter) ||
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
