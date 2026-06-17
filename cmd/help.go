package cmd

import "fmt"

func printHelp() {
	help := `Portman - Port Management CLI Tool

Usage:
  portman              Launch interactive TUI
  portman kill <port>  Kill the process listening on a port (direct, no prompt)
  portman list [port]  List listening ports
  portman upgrade      Upgrade to the latest version
  portman version      Show version information
  portman help         Show this help message

kill flags:
  -i, --interactive    Pick which process to kill when several listen on the port
      --all, -a        Also target non-listening connections (e.g. client sockets)
      --json           Output the result as JSON
  -q, --quiet          Print nothing on success (exit code only)

list flags:
      --all, -a        Include non-listening connections (default: listeners only)
      --json           Output as JSON

kill targets the listening (server) socket by default, so it won't kill a
client connection — e.g. a browser tab connected to your dev server.

Keybindings (TUI):
  ↑/↓ or j/k          Navigate
  Enter               Kill selected process
  r                   Refresh port list
  /                   Filter ports
  q or Ctrl+C         Quit

Examples:
  portman              # Launch interactive mode
  portman kill 3000    # Kill the server on port 3000
  portman list         # List listening ports
  portman list --json  # Machine-readable list

For scripts / AI agents:
  portman kill 3000          # one short line, no prompt: "✓ killed 3000 (node, pid 1234)"
  portman kill 3000 -q       # silent, rely on exit code
  portman kill 3000 --json   # structured result
  portman list 3000 --json   # just that port's listeners, as JSON
`
	fmt.Println(help)
}
