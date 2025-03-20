package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m model) View() string {
	switch m.mode {
	case modeFullDetail:
		title := lipgloss.NewStyle().Bold(true).Underline(true).Render("🔍 Full JSON Detail View")
		footer := lipgloss.NewStyle().Faint(true).Render("(↑↓ scroll, q/esc back)")

		// Compute visible lines
		start := m.detailOffset
		end := start + m.height - 4
		if end > len(m.fullDetailLines) {
			end = len(m.fullDetailLines)
		}
		content := strings.Join(m.fullDetailLines[start:end], "\n")

		return fmt.Sprintf("%s\n\n%s\n\n%s", title, content, footer)

	case modeRegexFilter:
		title := lipgloss.NewStyle().Bold(true).Underline(true).Render("🧹 Exclude Logs by Regex")
		helper := lipgloss.NewStyle().Faint(true).Render("(Enter = apply filter, Esc = cancel)")
		return title + "\n\n" + m.regexInput.View() + "\n" + helper
	case modePaste:
		title := lipgloss.NewStyle().Bold(true).Underline(true).Render("📋 Paste Mode")

		hint := lipgloss.NewStyle().Faint(true).Render(`
Instructions:
• Paste logs directly into the input area below (limited to ~99 lines).
• For larger log files, drag and drop a .log or .txt or .json file into the terminal to load automatically.
• Press Enter to continue or Esc to cancel.
`)
		return title + "\n\n" + m.textarea.View() + "\n" + hint + "\n\n(Enter = done, ctrl-z = clear, Esc = quit)"
	case modeView:
		var b strings.Builder

		filtered := m.pagedLogs()
		// Adjust these for your layout preference
		const messageStartColumn = 36

		for i, log := range filtered {
			globalIndex := m.offset + i
			prefix := "  "
			if i == m.cursor {
				prefix = "> "
			}
			indicator := "  "
			if len(log.Details) > 0 {
				if m.logs[globalIndex].Expanded {
					indicator = "⏷ " // down arrow = expanded
				} else {
					indicator = "⏵ " // right arrow = collapsed
				}
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
			line := indicator + header + spacing + log.Message

			// Render based on level
			switch level {
			case "ERROR", "WARN", "WARNING":
				b.WriteString(prefix + levelStyle.Render(line) + "\n")
			default:
				b.WriteString(prefix + white.Render(indicator+ts) + levelStyle.Render(lv) + white.Render(spacing+log.Message) + "\n")
			}

			if log.Expanded && len(log.Details) > 0 {
				b.WriteString(renderStyledJSON(log.Details) + "\n")
			}
		}

		if len(filtered) == 0 {
			b.WriteString(lipgloss.NewStyle().Faint(true).Render("No logs match the selected filter.\n"))
		}

		title := lipgloss.NewStyle().Bold(true).Underline(true).Render("📊 Log Viewer")
		helper := lipgloss.NewStyle().Faint(true).Render(
			title + "(q quit, z back, ↑↓ scroll, ⏎/space expand, e/w/i/d/a filter, r regex exclude, v view full JSON)",
		)
		b.WriteString("\n" + helper + "\n")

		return b.String()
	}
	return ""
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

	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		key := cyan.Render(fmt.Sprintf(`"%s"`, k))
		val := renderValue(data[k])
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
