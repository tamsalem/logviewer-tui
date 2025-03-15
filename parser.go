package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type logEntry struct {
	Level     string                 `json:"level"`
	Timestamp string                 `json:"timestamp"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"-"`
	Expanded  bool
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
