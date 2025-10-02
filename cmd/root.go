package cmd

import (
	"fmt"
	"os"

	"github.com/NoaTamburrini/portman/internal/tui"
	"github.com/NoaTamburrini/portman/internal/version"
)

// Execute is the main entry point for the CLI
func Execute() {
	// Check for updates in background (non-blocking, cached)
	version.CheckForUpdate()

	if len(os.Args) > 1 {
		// Handle subcommands
		switch os.Args[1] {
		case "kill":
			executeKill()
		case "help", "--help", "-h":
			printHelp()
		case "version", "--version", "-v":
			fmt.Printf("portman v%s\n", version.Version)
			os.Exit(0)
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
