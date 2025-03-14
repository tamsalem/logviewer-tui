# 📊 TUI Log Viewer

A terminal-based log viewer built with [BubbleTea](https://github.com/charmbracelet/bubbletea) and [LipGloss](https://github.com/charmbracelet/lipgloss). Paste your logs and explore them with keyboard navigation, color-coded levels, collapsible details, and filtering — all in your terminal.

---

## ✨ Features

- 📋 **Paste logs** directly into a terminal editor
- 🎨 **Color-coded levels**:
  - `ERROR` → red (full line)
  - `WARN` → yellow (full line)
  - `INFO` → blue (level only)
  - `DEBUG` → gray (level only)
- 📏 **Aligned log format** — messages always start at the same column
- 🔍 **Filter logs** by level: `e`, `w`, `i`, `d`, `a`
- 🔽 **Expand/collapse** remaining log fields (JSON object)
- 🧾 **Colorized JSON** for key/value pairs
- ⌨️ **Keyboard navigation** (see controls below)

---

## 📦 Installation

### 1. Clone the repository

```bash
git clone https://github.com/tamsalem/logviewer-tui.git
cd <your-repo-name>
```

### 2. Install dependencies

```bash
go mod tidy
```

This will fetch:
- `bubbletea`
- `bubbles/textarea`
- `lipgloss`

### 3. Run the app

```bash
go run main.go
```

---

## 🧑‍💻 How to Use

1. Paste your **JSON-formatted logs** (one per line)
2. Press `Enter` to enter viewer mode
3. Use the keyboard:

| Key        | Action                             |
|------------|------------------------------------|
| ↑ / ↓      | Move between log entries           |
| `Enter` / ␣ | Expand/collapse log details        |
| `e`        | Filter: show only `ERROR` logs     |
| `w`        | Filter: show only `WARN` logs      |
| `i`        | Filter: show only `INFO` logs      |
| `d`        | Filter: show only `DEBUG` logs     |
| `a`        | Show all logs                      |
| `q` / `Ctrl+C` | Quit the viewer                |

---

## 📎 Log Format Example

```json
{"level":"INFO","timestamp":"2025-03-13T16:05:36.013Z","message":"MongoDB initialized"}
{"level":"ERROR","timestamp":"2025-03-13T16:06:00.000Z","message":"Something failed","code":500}
```

Any additional fields will be available when expanding the log.

---

### ⚙️ Optional: Build an executable

```bash
go build -o logviewer
```

Run it directly:

```bash
./logviewer
```

Or move it to your system path:

```bash
sudo mv logviewer /usr/local/bin/
```

Now use it anywhere:

```bash
logviewer
```

---

## ✅ Requirements

- [Go](https://golang.org/doc/install) (v1.18+ recommended)
- Terminal that supports ANSI colors (iTerm2, Alacritty, VS Code terminal, etc.)

---

## 🛠 Roadmap / Ideas

- 🔍 Open for suggestions

---

## 📜 License

Apache-2.0 — use freely, build awesomely 🚀

---

## 💬 Credits

Built with:

- [BubbleTea](https://github.com/charmbracelet/bubbletea)
- [LipGloss](https://github.com/charmbracelet/lipgloss)
- [Bubbles](https://github.com/charmbracelet/bubbles)
