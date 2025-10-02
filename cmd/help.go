package cmd

import "fmt"

func printHelp() {
	help := `Portman - Port Management CLI Tool

Usage:
  portman              Launch interactive TUI
  portman kill <port>  Kill process on specific port
  portman help         Show this help message

Keybindings (TUI):
  ↑/↓ or j/k          Navigate
  Enter               Kill selected process
  r                   Refresh port list
  /                   Filter ports
  q or Ctrl+C         Quit

Examples:
  portman              # Launch interactive mode
  portman kill 3000    # Kill process on port 3000
`
	fmt.Println(help)
}
