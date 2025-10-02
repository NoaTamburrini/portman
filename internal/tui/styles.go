package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor   = lipgloss.Color("86")   // Cyan
	secondaryColor = lipgloss.Color("212")  // Pink
	successColor   = lipgloss.Color("42")   // Green
	errorColor     = lipgloss.Color("196")  // Red
	mutedColor     = lipgloss.Color("241")  // Gray
	selectedColor  = lipgloss.Color("219")  // Light purple

	// Title style
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	// Header row style
	headerStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(mutedColor)

	// Normal row style
	rowStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Selected row style
	selectedRowStyle = lipgloss.NewStyle().
				Foreground(selectedColor).
				Background(lipgloss.Color("235")).
				Bold(true).
				Padding(0, 1)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Padding(1, 0)

	// Status message styles
	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	// Filter input style
	filterStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	// Placeholder style
	placeholderStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Italic(true)
)
