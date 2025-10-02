package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var b strings.Builder

	// Title
	title := titleStyle.Render("🚢 PORTMAN - Port Manager")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Filter input (if in filter mode)
	if m.filterMode {
		b.WriteString(filterStyle.Render("Filter: "))
		b.WriteString(m.filterInput.View())
		b.WriteString("\n\n")
	}

	// Status message
	if m.statusMessage != "" {
		var statusStyle lipgloss.Style
		if m.statusIsError {
			statusStyle = errorStyle
		} else {
			statusStyle = successStyle
		}
		b.WriteString(statusStyle.Render(m.statusMessage))
		b.WriteString("\n\n")
	}

	// Port list
	if len(m.filteredPorts) == 0 {
		b.WriteString(renderMuted("No ports found"))
		b.WriteString("\n\n")
	} else {
		// Header
		header := fmt.Sprintf("%-8s %-10s %-8s %-20s %-30s",
			"PORT", "PROTOCOL", "PID", "PROCESS", "COMMAND")
		b.WriteString(headerStyle.Render(header))
		b.WriteString("\n")

		// Calculate how many rows we can show
		maxRows := m.height - 12 // Reserve space for title, status, help
		if maxRows < 5 {
			maxRows = 5
		}

		startIdx := 0
		endIdx := len(m.filteredPorts)

		// Adjust view if there are too many ports
		if len(m.filteredPorts) > maxRows {
			if m.cursor > maxRows/2 {
				startIdx = m.cursor - maxRows/2
			}
			if startIdx+maxRows > len(m.filteredPorts) {
				startIdx = len(m.filteredPorts) - maxRows
			}
			endIdx = startIdx + maxRows
		}

		// Rows
		for i := startIdx; i < endIdx; i++ {
			p := m.filteredPorts[i]

			// Truncate command if too long
			command := p.Command
			if len(command) > 30 {
				command = command[:27] + "..."
			}

			row := fmt.Sprintf("%-8d %-10s %-8d %-20s %-30s",
				p.Number,
				p.Protocol,
				p.PID,
				truncate(p.ProcessName, 20),
				command,
			)

			// Apply style based on selection
			if i == m.cursor {
				row = "▸ " + row
				b.WriteString(selectedRowStyle.Render(row))
			} else {
				row = "  " + row
				b.WriteString(rowStyle.Render(row))
			}
			b.WriteString("\n")
		}

		// Show scroll indicator if needed
		if len(m.filteredPorts) > maxRows {
			indicator := fmt.Sprintf("  [Showing %d-%d of %d]", startIdx+1, endIdx, len(m.filteredPorts))
			b.WriteString(renderMuted(indicator))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")

	// Help text
	help := ""
	if m.filterMode {
		help = "Enter: apply filter • Esc: cancel"
	} else if m.confirmingKill {
		help = "y: confirm kill • n: cancel"
	} else {
		help = "↑/↓ j/k: navigate • Enter: kill • r: refresh • /: filter • q: quit"
	}
	b.WriteString(helpStyle.Render(help))

	return b.String()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func renderMuted(s string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(s)
}
