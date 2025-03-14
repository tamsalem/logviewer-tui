package main

import (
	"encoding/json"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

type mode int

const (
	modePaste mode = iota
	modeView
)

type logEntry struct {
	Level     string                 `json:"level"`
	Timestamp string                 `json:"timestamp"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"-"`
	Expanded  bool
}

type model struct {
	mode     mode
	cursor   int
	filter   string
	logs     []logEntry
	textarea textarea.Model
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Paste logs here and press Enter when done..."
	ta.Focus()
	ta.CharLimit = 0
	ta.SetHeight(20)

	return model{
		mode:     modePaste,
		textarea: ta,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case modePaste:
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				input := m.textarea.Value()
				parsed := parseLogs(input)
				if len(parsed) == 0 {
					m.textarea.Placeholder = "⚠️ No valid logs found. Try again."
					m.textarea.SetValue("")
					return m, nil
				}
				m.logs = parsed
				m.mode = modeView
				return m, nil
			case "ctrl+c", "esc":
				return m, tea.Quit
			}
		}
		return m, cmd

	case modeView:
		switch key := msg.(type) {
		case tea.KeyMsg:
			switch key.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down":
				if m.cursor < len(m.filteredLogs())-1 {
					m.cursor++
				}
			case "enter", " ":
				logs := m.filteredLogs()
				if len(logs) > m.cursor {
					i := m.findLogIndex(logs[m.cursor])
					m.logs[i].Expanded = !m.logs[i].Expanded
				}
			case "e", "w", "i", "d", "a":
				switch key.String() {
				case "e":
					m.filter = "ERROR"
				case "w":
					m.filter = "WARN"
				case "i":
					m.filter = "INFO"
				case "d":
					m.filter = "DEBUG"
				case "a":
					m.filter = ""
				}
				m.cursor = 0
			}
		}
	}

	return m, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}


func (m model) View() string {
	switch m.mode {
	case modePaste:
		title := lipgloss.NewStyle().Bold(true).Underline(true).Render("📋 Paste Mode")
		return title + "\n\n" + m.textarea.View() + "\n\n(Enter = done, Esc = quit)"
	case modeView:
		var b strings.Builder
		title := lipgloss.NewStyle().Bold(true).Underline(true).Render("📊 Log Viewer")
		b.WriteString(title + " (q to quit, arrows to move, enter/space to expand, e/w/i/d/a to filter)\n\n")

		filtered := m.filteredLogs()
		// Adjust these for your layout preference
		const messageStartColumn = 36

		for i, log := range filtered {
			prefix := "  "
			if i == m.cursor {
				prefix = "> "
			}

			level := strings.ToUpper(log.Level)
			levelStyle := levelColor(level)
			white := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

			// Format core line parts
			ts := fmt.Sprintf("[%s]", log.Timestamp)
			lv := fmt.Sprintf("[%s]", level)
			header := ts + lv

			// Compute padding so message always starts at column N
			headerWidth := lipgloss.Width(prefix + header)
			spacing := strings.Repeat(" ", max(0, messageStartColumn-headerWidth))
			line := header + spacing + log.Message

			// Render based on level
			switch level {
			case "ERROR", "WARN", "WARNING":
				b.WriteString(prefix + levelStyle.Render(line) + "\n")
			default:
				b.WriteString(prefix + white.Render(ts) + levelStyle.Render(lv) + white.Render(spacing+log.Message) + "\n")
			}

			if log.Expanded && len(log.Details) > 0 {
				b.WriteString(renderStyledJSON(log.Details) + "\n")
			}
		}

		if len(filtered) == 0 {
			b.WriteString(lipgloss.NewStyle().Faint(true).Render("No logs match the selected filter.\n"))
		}

		return b.String()
	}
	return ""
}

func (m model) filteredLogs() []logEntry {
	if m.filter == "" {
		return m.logs
	}
	var filtered []logEntry
	for _, log := range m.logs {
		if strings.EqualFold(log.Level, m.filter) {
			filtered = append(filtered, log)
		}
	}
	return filtered
}

func (m model) findLogIndex(target logEntry) int {
	for i, log := range m.logs {
		if log.Timestamp == target.Timestamp && log.Message == target.Message && log.Level == target.Level {
			return i
		}
	}
	return 0
}

func parseLogs(input string) []logEntry {
	lines := strings.Split(input, "\n")
	var logs []logEntry

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var raw map[string]interface{}
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			continue
		}

		log := logEntry{
			Level:     fmt.Sprintf("%v", raw["level"]),
			Timestamp: fmt.Sprintf("%v", raw["timestamp"]),
			Message:   fmt.Sprintf("%v", raw["message"]),
			Details:   make(map[string]interface{}),
		}

		for k, v := range raw {
			switch k {
			case "level", "timestamp", "message", "jobName", "traceId", "requestId", "workflowId", "currentExecutedFlow":
				continue
			default:
				log.Details[k] = v
			}
		}

		logs = append(logs, log)
	}
	return logs
}

func renderStyledJSON(data map[string]interface{}) string {
	var b strings.Builder
	cyan := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))    // keys
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))   // strings
	yellow := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))  // numbers
	magenta := lipgloss.NewStyle().Foreground(lipgloss.Color("5")) // booleans
	gray := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))    // null

	renderValue := func(v interface{}) string {
		switch val := v.(type) {
		case string:
			return green.Render(fmt.Sprintf(`"%s"`, val))
		case float64, int:
			return yellow.Render(fmt.Sprintf("%v", val))
		case bool:
			return magenta.Render(fmt.Sprintf("%v", val))
		case nil:
			return gray.Render("null")
		default:
			encoded, _ := json.Marshal(val)
			return green.Render(string(encoded))
		}
	}

	for k, v := range data {
		key := cyan.Render(fmt.Sprintf(`"%s"`, k))
		val := renderValue(v)
		b.WriteString(fmt.Sprintf("    %s: %s\n", key, val))
	}

	return b.String()
}

func levelColor(level string) lipgloss.Style {
	switch strings.ToUpper(level) {
	case "ERROR":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // Red
	case "WARN", "WARNING":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // Yellow
	case "INFO":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("4")) // Blue
	case "DEBUG":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("8")) // Gray
	default:
		return lipgloss.NewStyle()
	}
}

func main() {
	if err := tea.NewProgram(initialModel()).Start(); err != nil {
		fmt.Println("Error:", err)
	}
}
