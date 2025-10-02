package tui

import (
	"fmt"

	"github.com/NoaTamburrini/portman/internal/process"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Handle filter mode separately
		if m.filterMode {
			return m.handleFilterMode(msg)
		}

		// Handle confirmation mode
		if m.confirmingKill {
			return m.handleConfirmMode(msg)
		}

		// Normal mode key handling
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.filteredPorts)-1 {
				m.cursor++
			}

		case "r":
			m.scanning = true
			m.statusMessage = "Refreshing..."
			m.statusIsError = false
			return m, scanPorts

		case "/":
			m.filterMode = true
			m.filterInput.Focus()
			return m, textinput.Blink

		case "enter":
			if len(m.filteredPorts) > 0 {
				m.confirmingKill = true
				selectedPort := m.filteredPorts[m.cursor]
				m.statusMessage = fmt.Sprintf("Kill process on port %d (PID: %d)? [y/N]",
					selectedPort.Number, selectedPort.PID)
				m.statusIsError = false
			}
		}

	case scanCompleteMsg:
		m.scanning = false
		if msg.err != nil {
			m.statusMessage = fmt.Sprintf("Error: %v", msg.err)
			m.statusIsError = true
		} else {
			m.ports = msg.ports
			m.filterPorts()
			m.statusMessage = fmt.Sprintf("Found %d active port(s)", len(m.ports))
			m.statusIsError = false
		}

	case killCompleteMsg:
		if msg.success {
			m.statusMessage = fmt.Sprintf("✓ %s", msg.message)
			m.statusIsError = false
			// Refresh after kill
			return m, scanPorts
		} else {
			m.statusMessage = fmt.Sprintf("✗ %s", msg.message)
			m.statusIsError = true
		}
	}

	return m, nil
}

func (m Model) handleFilterMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "ctrl+c":
		m.filterMode = false
		m.filterInput.Blur()
		m.filterInput.SetValue("")
		m.filterPorts()
		return m, nil

	case "enter":
		m.filterMode = false
		m.filterInput.Blur()
		m.filterPorts()
		return m, nil
	}

	var cmd tea.Cmd
	m.filterInput, cmd = m.filterInput.Update(msg)
	m.filterPorts()
	return m, cmd
}

func (m Model) handleConfirmMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.confirmingKill = false

	switch msg.String() {
	case "y", "Y":
		if len(m.filteredPorts) > 0 {
			selectedPort := m.filteredPorts[m.cursor]
			m.statusMessage = fmt.Sprintf("Killing process on port %d...", selectedPort.Number)
			m.statusIsError = false

			// Kill the process
			return m, func() tea.Msg {
				result := process.KillProcess(selectedPort.PID)
				return killCompleteMsg{
					success: result.Success,
					message: result.Message,
				}
			}
		}
	}

	m.statusMessage = "Kill cancelled"
	m.statusIsError = false
	return m, nil
}
