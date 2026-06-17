# 🚢 Portman

A fast and intuitive CLI tool for managing ports and processes on your system. Kill processes by port number with an interactive TUI or simple commands.

## ✨ Features

- 🖥️ **Interactive TUI** - Beautiful terminal interface for browsing active ports
- ⚡ **Quick Kill** - Instantly kill processes by port number
- 🔍 **Smart Filtering** - Search and filter ports in real-time
- 🎯 **Simple Commands** - Easy-to-use CLI for automation
- 💻 **Cross-Platform** - Works on macOS, Linux, and Windows

## 📦 Installation

### Quick Install (Recommended)

**macOS / Linux:**

```bash
curl -fsSL https://raw.githubusercontent.com/NoaTamburrini/portman/main/install.sh | sh
```

**Windows (PowerShell):**

```powershell
irm https://raw.githubusercontent.com/NoaTamburrini/portman/main/install.ps1 | iex
```

These scripts automatically detect your architecture and install the latest version.

### Manual Installation

<details>
<summary>Click to expand manual installation instructions</summary>

#### macOS (Apple Silicon)
```bash
curl -L https://github.com/NoaTamburrini/portman/releases/download/v1.0.0/portman-darwin-arm64 -o /usr/local/bin/portman && chmod +x /usr/local/bin/portman
```

#### macOS (Intel)
```bash
curl -L https://github.com/NoaTamburrini/portman/releases/download/v1.0.0/portman-darwin-amd64 -o /usr/local/bin/portman && chmod +x /usr/local/bin/portman
```

#### Linux
```bash
curl -L https://github.com/NoaTamburrini/portman/releases/download/v1.0.0/portman-linux-amd64 -o /usr/local/bin/portman && chmod +x /usr/local/bin/portman
```

#### Windows
Download the latest `.zip` from the [releases page](https://github.com/NoaTamburrini/portman/releases), extract `portman.exe`, and add it to your `PATH`.

</details>

## 🚀 Usage

### Interactive Mode

Launch the TUI to browse and kill processes:

```bash
portman
```

**Keybindings:**

- `↑/↓` or `j/k` - Navigate through ports
- `Enter` - Kill selected process
- `r` - Refresh port list
- `/` - Filter/search ports
- `q` or `Ctrl+C` - Quit

### Command Mode

Kill the process listening on a port. This is **direct** — no prompt, no menu:

```bash
portman kill 3000
# ✓ killed 3000 (node, pid 1234)
```

`kill` targets the **listening (server)** socket by default, so it won't kill a
client connection that happens to share the port — for example a browser tab
connected to your dev server. Pass `--all` if you really want to kill those too.

List the listening ports:

```bash
portman list
```

### Help

```bash
portman help
```

## 🛠️ Examples

```bash
# Launch interactive TUI
portman

# Kill process running on port 3000
portman kill 3000

# Show help
portman --help
```

## 🤖 Scripting / AI agent usage

Every command works non-interactively, so portman is safe to drive from scripts
or AI agents without anything blocking on a terminal prompt.

```bash
# Kill instantly — one short line, no prompt, no menu
portman kill 3000              # ✓ killed 3000 (node, pid 1234)

# Silent: print nothing, rely on the exit code
portman kill 3000 -q

# Structured output for parsing
portman kill 3000 --json       # [{"port":3000,"pid":1234,"process":"node","success":true,"message":"..."}]

# List only listening ports (compact), as JSON
portman list --json

# Just one port's listeners
portman list 3000 --json
```

Key points for agents:

- `kill` only targets the **LISTEN**ing process by default — it never kills a
  client connection (e.g. your browser) that shares the port number.
- Multiple listeners on one port (e.g. IPv4 + IPv6) are all killed directly. The
  interactive picker is opt-in via `-i` and is only for human use.
- `list` shows **listeners only** by default, keeping output small. Use `--all`
  for the full picture including established connections.

## 📝 License

MIT

## 🤝 Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

---

Made with ❤️ by [NoaTamburrini](https://github.com/NoaTamburrini)
