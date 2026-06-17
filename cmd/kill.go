package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/NoaTamburrini/portman/internal/process"
	"github.com/NoaTamburrini/portman/internal/scanner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type killOptions struct {
	interactive bool
	jsonOutput  bool
	quiet       bool
	all         bool
}

// killResult is the per-process outcome emitted with --json.
type killResult struct {
	Port    int    `json:"port"`
	PID     int    `json:"pid"`
	Process string `json:"process"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func executeKill() {
	portStr := ""
	opts := killOptions{}
	for _, arg := range os.Args[2:] {
		switch arg {
		case "-i", "--interactive":
			opts.interactive = true
		case "--json":
			opts.jsonOutput = true
		case "-q", "--quiet":
			opts.quiet = true
		case "--all", "-a":
			opts.all = true
		default:
			if portStr == "" {
				portStr = arg
			}
		}
	}

	if portStr == "" {
		fmt.Fprintln(os.Stderr, "Usage: portman kill <port> [-i] [--all] [--json] [-q]")
		os.Exit(1)
	}

	portNum, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid port number: %s\n", portStr)
		os.Exit(1)
	}

	if portNum < 1 || portNum > 65535 {
		fmt.Fprintf(os.Stderr, "Port number must be between 1 and 65535\n")
		os.Exit(1)
	}

	ports, err := scanner.ScanPorts()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning ports: %v\n", err)
		os.Exit(1)
	}

	// Pick targets: listeners only by default, every match with --all.
	var targets []scanner.Port
	if opts.all {
		targets = scanner.FindAllByPort(ports, portNum)
	} else {
		targets = scanner.FindKillTargets(ports, portNum)
	}

	if len(targets) == 0 {
		if opts.jsonOutput {
			fmt.Println("[]")
		} else if !opts.quiet {
			if len(scanner.FindAllByPort(ports, portNum)) > 0 {
				fmt.Fprintf(os.Stderr, "No process listening on port %d (use --all to kill connections)\n", portNum)
			} else {
				fmt.Fprintf(os.Stderr, "No process found on port %d\n", portNum)
			}
		}
		os.Exit(1)
	}

	// Interactive picker is opt-in only; it never runs for scripts/agents.
	if opts.interactive && len(targets) > 1 {
		selected := showSelectionMenu(targets, portNum)
		if selected == nil {
			if !opts.quiet {
				fmt.Println("Cancelled")
			}
			return
		}
		targets = selected
	}

	killTargets(targets, opts)
}

// killTargets kills each target and reports results per the output options.
func killTargets(targets []scanner.Port, opts killOptions) {
	results := make([]killResult, 0, len(targets))
	anyFailed := false

	for _, p := range targets {
		res := process.KillProcess(p.PID)
		if !res.Success {
			anyFailed = true
		}
		results = append(results, killResult{
			Port:    p.Number,
			PID:     p.PID,
			Process: p.ProcessName,
			Success: res.Success,
			Message: res.Message,
		})
	}

	switch {
	case opts.jsonOutput:
		out, _ := json.Marshal(results)
		fmt.Println(string(out))
	case opts.quiet:
		// exit code only
	default:
		for _, r := range results {
			if r.Success {
				fmt.Printf("✓ killed %d (%s, pid %d)\n", r.Port, r.Process, r.PID)
			} else {
				fmt.Fprintf(os.Stderr, "✗ port %d (%s, pid %d): %s\n", r.Port, r.Process, r.PID, r.Message)
			}
		}
	}

	if anyFailed {
		os.Exit(1)
	}
}

type selectionModel struct {
	choices   []scanner.Port
	cursor    int
	selected  map[int]bool
	portNum   int
	quitting  bool
	cancelled bool
}

func (m selectionModel) Init() tea.Cmd {
	return nil
}

func (m selectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.cancelled = true
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)+1 {
				m.cursor++
			}
		case "enter":
			if m.cursor == len(m.choices) {
				// "Kill all selected" - just quit with current selection
				m.quitting = true
				return m, tea.Quit
			} else if m.cursor == len(m.choices)+1 {
				// "Cancel" option
				m.cancelled = true
				m.quitting = true
				return m, tea.Quit
			} else {
				// Toggle individual selection
				m.selected[m.cursor] = !m.selected[m.cursor]
			}
		case " ":
			// Space to toggle (only for process items)
			if m.cursor < len(m.choices) {
				m.selected[m.cursor] = !m.selected[m.cursor]
			}
		}
	}
	return m, nil
}

func (m selectionModel) View() string {
	// Styles matching main TUI
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true).Padding(0, 1)
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true).
		BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(lipgloss.Color("241"))
	selectedRowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("219")).
		Background(lipgloss.Color("235")).Bold(true).Padding(0, 1)
	rowStyle := lipgloss.NewStyle().Padding(0, 1)
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Padding(1, 0)
	checkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)

	var s string

	// Title
	s += titleStyle.Render(fmt.Sprintf("🔍 Select processes to kill on port %d", m.portNum)) + "\n\n"

	// Header
	header := fmt.Sprintf("%-8s %-10s %-8s %-20s", "SELECT", "PORT", "PID", "PROCESS")
	s += headerStyle.Render(header) + "\n"

	// Process rows
	for i, choice := range m.choices {
		checkbox := "[ ]"
		if m.selected[i] {
			checkbox = checkStyle.Render("[✓]")
		}

		row := fmt.Sprintf("%-8s %-10d %-8d %-20s",
			checkbox,
			choice.Number,
			choice.PID,
			choice.ProcessName)

		if m.cursor == i {
			s += selectedRowStyle.Render("▸ "+row) + "\n"
		} else {
			s += rowStyle.Render("  "+row) + "\n"
		}
	}

	s += "\n"

	// Action options
	killAllRow := "Kill all selected"
	if m.cursor == len(m.choices) {
		s += selectedRowStyle.Render("▸ "+killAllRow) + "\n"
	} else {
		s += rowStyle.Render("  "+killAllRow) + "\n"
	}

	cancelRow := "Cancel"
	if m.cursor == len(m.choices)+1 {
		s += selectedRowStyle.Render("▸ "+cancelRow) + "\n"
	} else {
		s += rowStyle.Render("  "+cancelRow) + "\n"
	}

	s += "\n"

	// Help
	help := "↑/↓ j/k: navigate • Space: toggle • Enter: confirm • q/Esc: cancel"
	s += helpStyle.Render(help)

	return s
}

func showSelectionMenu(matches []scanner.Port, portNum int) []scanner.Port {
	m := selectionModel{
		choices:  matches,
		selected: make(map[int]bool),
		portNum:  portNum,
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	result := finalModel.(selectionModel)

	if result.cancelled {
		return nil
	}

	var selected []scanner.Port
	for i, port := range result.choices {
		if result.selected[i] {
			selected = append(selected, port)
		}
	}

	return selected
}
