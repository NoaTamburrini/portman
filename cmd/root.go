package cmd

import (
	"fmt"
	"os"

	"github.com/NoaTamburrini/portman/internal/tui"
	"github.com/NoaTamburrini/portman/internal/version"
)

// Execute is the main entry point for the CLI
func Execute() {
	if len(os.Args) > 1 {
		// Check for updates only in command mode (TUI handles it internally)
		version.CheckForUpdate()
		// Handle subcommands
		switch os.Args[1] {
		case "kill":
			executeKill()
		case "list", "ls":
			executeList()
		case "upgrade":
			executeUpgrade()
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
