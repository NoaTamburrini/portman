package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/NoaTamburrini/portman/internal/process"
	"github.com/NoaTamburrini/portman/internal/scanner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	// Scan to find all processes on the port
	ports, err := scanner.ScanPorts()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning ports: %v\n", err)
		os.Exit(1)
	}

	matches := scanner.FindAllByPort(ports, portNum)
	if len(matches) == 0 {
		fmt.Printf("No process found on port %d\n", portNum)
		os.Exit(1)
	}

	// If only one process, kill it directly
	if len(matches) == 1 {
		port := matches[0]
		fmt.Printf("Killing process on port %d (PID: %d, Process: %s)...\n",
			port.Number, port.PID, port.ProcessName)

		result := process.KillProcess(port.PID)
		if result.Success {
			fmt.Printf("✓ %s\n", result.Message)
		} else {
			fmt.Fprintf(os.Stderr, "✗ %s\n", result.Message)
			os.Exit(1)
		}
		return
	}

	// Multiple processes - show Bubble Tea selection menu
	selected := showSelectionMenu(matches, portNum)

	if selected == nil {
		fmt.Println("Cancelled")
		return
	}

	// Kill all selected
	if len(selected) == len(matches) {
		fmt.Printf("Killing all %d processes on port %d...\n", len(selected), portNum)
	}

	for _, p := range selected {
		fmt.Printf("Killing PID %d (%s)...\n", p.PID, p.ProcessName)
		result := process.KillProcess(p.PID)
		if result.Success {
			fmt.Printf("✓ %s\n", result.Message)
		} else {
			fmt.Fprintf(os.Stderr, "✗ %s\n", result.Message)
		}
	}
}

type selectionModel struct {
	choices  []scanner.Port
	cursor   int
	selected map[int]bool
	portNum  int
	quitting bool
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
