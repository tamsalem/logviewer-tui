package main

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	modePaste int = iota
	modeView
	modeRegexFilter // 👈 new mode
)

type model struct {
	mode            int
	cursor          int
	filter          string
	logs            []logEntry
	textarea        textarea.Model
	height          int
	offset          int
	regexInput      textarea.Model
	excludePatterns []*regexp.Regexp
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Paste logs here and press Enter when done..."
	ta.Focus()
	ta.CharLimit = 0
	ta.SetHeight(20)
	regexTA := textarea.New()
	regexTA.Placeholder = "Comma-separated regex to exclude (e.g. *debug*,*heartbeat$)"
	regexTA.CharLimit = 0
	regexTA.SetHeight(3)

	return model{
		mode:       modePaste,
		textarea:   ta,
		regexInput: regexTA,
	}
}

func (m model) pageSize() int {
	size := m.height - 6
	if size < 1 {
		return 1
	}
	return size
}

func (m *model) scrollUp() {
	if m.cursor > 0 {
		m.cursor--
	} else if m.offset > 0 {
		m.offset--
	}
}

func (m *model) scrollDown() {
	if m.cursor < m.pageSize()-1 && m.cursor < len(m.pagedLogs())-1 {
		m.cursor++
	} else if m.offset+m.cursor+1 < len(m.filteredLogs()) {
		m.offset++
	}
}

func (m model) findLogIndex(target logEntry) int {
	for i, log := range m.logs {
		if log.Timestamp == target.Timestamp && log.Message == target.Message && log.Level == target.Level {
			return i
		}
	}
	return 0
}

func (m model) filteredLogs() []logEntry {
	var filtered []logEntry
outer:
	for _, log := range m.logs {
		if m.filter != "" && !strings.EqualFold(log.Level, m.filter) {
			continue
		}
		combined := log.Message + " " + log.Level + " " + log.Timestamp
		for _, re := range m.excludePatterns {
			if re.MatchString(combined) {
				continue outer
			}
		}
		filtered = append(filtered, log)
	}
	return filtered
}

func (m model) pagedLogs() []logEntry {
	var visible []logEntry
	logs := m.filteredLogs()
	linesAvailable := m.height - 6 // adjust for header/footer

	// Start from offset
	for i := m.offset; i < len(logs); i++ {
		log := logs[i]
		linesUsed := 1
		if log.Expanded && len(log.Details) > 0 {
			linesUsed += strings.Count(renderStyledJSON(log.Details), "\n")
		}

		if linesAvailable-linesUsed < 0 {
			break
		}

		visible = append(visible, log)
		linesAvailable -= linesUsed
	}

	return visible
}

func compileRegexList(input string) []*regexp.Regexp {
	var patterns []*regexp.Regexp
	for _, s := range strings.Split(input, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if re, err := regexp.Compile(s); err == nil {
			patterns = append(patterns, re)
		}
	}
	return patterns
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
	}
	switch m.mode {
	case modeRegexFilter:
		var cmd tea.Cmd
		m.regexInput, cmd = m.regexInput.Update(msg)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				input := m.regexInput.Value()
				m.excludePatterns = compileRegexList(input)
				m.mode = modeView
				m.cursor = 0
				m.offset = 0
				return m, nil
			case "esc", "ctrl+c":
				m.mode = modeView
				return m, nil
			}
		}
		return m, cmd

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
				m.scrollUp()
			case "down":
				m.scrollDown()
			case "enter", " ":
				logs := m.pagedLogs()
				if len(logs) > m.cursor {
					i := m.findLogIndex(logs[m.cursor])
					m.logs[i].Expanded = !m.logs[i].Expanded
				}
			case "e", "w", "i", "d", "a", "r":
				m.cursor = 0
				m.offset = 0

				// Reset all expanded states
				for i := range m.logs {
					m.logs[i].Expanded = false
				}

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
					m.excludePatterns = nil
				case "r":
					m.mode = modeRegexFilter
					m.regexInput.Focus()
					m.regexInput.SetValue("")
				}
			}
		}
	}

	return m, nil
}
