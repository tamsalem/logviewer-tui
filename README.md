# 📊 TUI Log Viewer

A terminal-based log viewer built with [BubbleTea](https://github.com/charmbracelet/bubbletea) and [LipGloss](https://github.com/charmbracelet/lipgloss). Easily explore structured logs with keyboard navigation, color-coded levels, collapsible metadata, filtering, and more — all in your terminal.

---

## ✨ Features

- 📋 Paste JSON logs directly into the terminal
- 🎨 Color-coded log levels:
  - `ERROR` → Red (entire line)
  - `WARN`  → Yellow (entire line)
  - `INFO`  → Blue (level only)
  - `DEBUG` → Gray (level only)
- 📐 Aligned log formatting — message text always starts at the same column
- 🔍 Level-based filtering: `e`, `w`, `i`, `d`, `a`
- 🔽 Expand/collapse log fields with pretty-printed, colorized JSON
- 🧠 Regex-based exclusion filtering
- ⌨️ Keyboard-first navigation for log review
- ⚙️ Support for fetching logs from Argo Workflows (`--workflow` flag)

---

## ✅ Requirements

- [Go](https://golang.org/doc/install) 1.18 or higher
- A terminal that supports ANSI colors (e.g. iTerm2, Alacritty, VS Code)

---

## 📦 Installation

### 1. Clone the repository

```bash
git clone https://github.com/tamsalem/logviewer-tui.git
cd logviewer-tui
```

### 2. Install dependencies

```bash
go mod tidy
```

This pulls all required libraries:
- `bubbletea`
- `bubbles/textarea`
- `lipgloss`

### 3. Run

```bash
go run .
```

### 4. (Optional) Build the executable

```bash
go build -o logviewer
```

### 5. (Optional) Move it to your `$PATH`

```bash
sudo mv logviewer /usr/local/bin/
```

Now you can run it globally:

```bash
logviewer
```

---

## 🚀 Usage

```bash
logviewer
```

1. Paste JSON-formatted logs (one object per line)
2. Press `Enter` to parse and enter viewer mode
3. Navigate using the keyboard

---

## 🏷️ Command-Line Flags

| Flag           | Description                                                          |
|----------------|----------------------------------------------------------------------|
| `--workflow`   | (Optional) Provide an Argo Workflow name to fetch logs from Argo API |

### Example

```bash
logviewer --workflow 9f9aab90-319b-4655-905c-7ea2db0ef550
```

- Connects to your local Argo server (`http://localhost:2746`)
- Prompts you to select a workflow step
- Loads and renders the logs for that step

---

## ⌨️ Controls

| Key                | Action                                           |
|--------------------|--------------------------------------------------|
| ↑ / ↓              | Navigate between log entries                     |
| `Enter` / ␣        | Expand or collapse log metadata                  |
| `e`                | Filter: only show `ERROR` logs                   |
| `w`                | Filter: only show `WARN` logs                    |
| `i`                | Filter: only show `INFO` logs                    |
| `d`                | Filter: only show `DEBUG` logs                   |
| `a`                | Reset filters and show all logs                  |
| `r`                | Set regex to exclude logs (comma-separated)      |
| `v`                | View full details (pretty JSON) in full-screen   |
| `home/end` / `g/G`  | Jump to top / bottom                             |
| `q` / Ctrl+C       | Quit the viewer                                  |

---

## 💡 Paste Mode Tips

- For large logs, use **drag-and-drop** to insert a `.json` or `.log` file into the terminal
- Paste mode supports up to ~99 lines directly
- Logs must be line-delimited JSON objects

---

## 🧾 Example Log Format

```json
{"level":"INFO","timestamp":"2025-03-13T16:05:36.013Z","message":"MongoDB initialized"}
{"level":"ERROR","timestamp":"2025-03-13T16:06:00.000Z","message":"Something failed","code":500}
```

Any additional fields (e.g. `code`, `context`) will be shown when expanded.

---

## 📜 License

Apache-2.0 — use freely, build awesomely 🚀

---

## 🙌 Credits

Powered by:

- [BubbleTea](https://github.com/charmbracelet/bubbletea)
- [LipGloss](https://github.com/charmbracelet/lipgloss)
- [Bubbles](https://github.com/charmbracelet/bubbles)
