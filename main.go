package main

import (
	"flag"
	"fmt"

	"logviewer-tui/argo"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	workflow := flag.String("workflow", "", "Argo workflow name e.g 'sync-customer-template-gj97n'")
	flag.Parse()
	m := initialModel()
	if *workflow != "" {
		logs, err := argo.RunWorkflowMode(*workflow)
		if err != nil || logs == "" {
			fmt.Println("❌ Failed to fetch logs:", err)
			return
		}
		m.logs = parseLogs(logs)
		m.mode = modeView
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
