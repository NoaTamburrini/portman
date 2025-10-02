package cmd

import (
	"fmt"
	"os"

	"github.com/NoaTamburrini/portman/internal/tui"
)

// Execute is the main entry point for the CLI
func Execute() {
	if len(os.Args) > 1 {
		// Handle subcommands
		switch os.Args[1] {
		case "kill":
			executeKill()
		case "help", "--help", "-h":
			printHelp()
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			printHelp()
			os.Exit(1)
		}
	} else {
		// No arguments - launch TUI
		tui.Start()
	}
}
