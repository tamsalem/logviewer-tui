package argo

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item string

func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }

type listModel struct {
	list     list.Model
	selected string
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected := m.list.SelectedItem()
			if selected != nil {
				m.selected = selected.FilterValue()
				return m, tea.Quit
			}
		case "q", "esc":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listModel) View() string {
	return m.list.View()
}

func promptStepSelection(steps []string) (string, error) {
	items := make([]list.Item, len(steps))
	for i, step := range steps {
		items[i] = item(step)
	}

	l := list.New(items, list.NewDefaultDelegate(), 50, 20)
	l.Title = "Select a step to view logs"

	m := listModel{list: l}
	program := tea.NewProgram(m)
	ListModel, err := program.Run()
	if err != nil {
		return "", err
	}
	// Type assert the final model back to listModel
	if lm, ok := ListModel.(listModel); ok {
		fmt.Println("✅ Step selected:", lm.selected)
		return lm.selected, nil
	}

	return "", fmt.Errorf("failed to cast final model")
}

func RunWorkflowMode(workflow string) (string, error) {
	token, err := getArgoToken()
	if err != nil {
		fmt.Println("❌ Failed to get token:", err)
		return "", err
	}

	workflowUid, steps, idMap, err := fetchWorkflowSteps(workflow, token)
	if err != nil {
		fmt.Println("❌ Failed to fetch workflow:", err)
		return "", err
	}

	step, err := promptStepSelection(steps)
	if err != nil {
		fmt.Println("❌ Failed to select step:", err)
		return "", err
	}

	if step == "" {
		fmt.Println("⚠️ Step selection failed — no name returned.")
	}

	nodeID := idMap[step]
	logs, err := fetchLogs(workflow, workflowUid, nodeID, token)
	if err != nil {
		fmt.Println("❌ Failed to fetch logs:", err)
		return "", err
	}
	return logs, nil
}
