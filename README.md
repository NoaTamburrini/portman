# 🚢 Portman

A fast and intuitive CLI tool for managing ports and processes on your system. Kill processes by port number with an interactive TUI or simple commands.

## ✨ Features

- 🖥️ **Interactive TUI** - Beautiful terminal interface for browsing active ports
- ⚡ **Quick Kill** - Instantly kill processes by port number
- 🔍 **Smart Filtering** - Search and filter ports in real-time
- 🎯 **Simple Commands** - Easy-to-use CLI for automation
- 💻 **Cross-Platform** - Works on macOS and Linux

## 📦 Installation

### macOS (Apple Silicon)
```bash
curl -L https://github.com/NoaTamburrini/portman/releases/download/v1.0.0/portman-darwin-arm64 -o /usr/local/bin/portman && chmod +x /usr/local/bin/portman
```

### macOS (Intel)
```bash
curl -L https://github.com/NoaTamburrini/portman/releases/download/v1.0.0/portman-darwin-amd64 -o /usr/local/bin/portman && chmod +x /usr/local/bin/portman
```

### Linux
```bash
curl -L https://github.com/NoaTamburrini/portman/releases/download/v1.0.0/portman-linux-amd64 -o /usr/local/bin/portman && chmod +x /usr/local/bin/portman
```

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
Kill a process on a specific port:
```bash
portman kill 3000
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

## 📝 License

MIT

## 🤝 Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

---

Made with ❤️ by [NoaTamburrini](https://github.com/NoaTamburrini)
