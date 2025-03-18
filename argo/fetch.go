package argo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
)

func sendLogRequest(client *http.Client, url, token string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Argo API error (%d): %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

func fetchLogs(workflowName, workflowUid, nodeID, token string) (string, error) {
	if nodeID == "" {
		return "", fmt.Errorf("nodeID is empty")
	}

	client := &http.Client{}

	// Primary logs URL
	primaryURL := fmt.Sprintf("http://localhost:2746/artifact-files/cas/workflows/%s/%s/outputs/main-logs", workflowName, nodeID)
	fmt.Println("📡 Fetching logs from:", primaryURL)

	primaryResp, err := sendLogRequest(client, primaryURL, token)
	if err == nil {
		return primaryResp, nil
	}

	fmt.Println("⚠️ Primary log fetch failed, trying fallback (archived-workflows)")

	// Fallback logs URL
	fallbackURL := fmt.Sprintf("http://localhost:2746/artifact-files/cas/archived-workflows/%s/%s/outputs/main-logs", workflowUid, nodeID)
	fmt.Println("📡 Fetching logs from fallback:", fallbackURL)

	fallbackResp, err := sendLogRequest(client, fallbackURL, token)
	if err != nil {
		return "", fmt.Errorf("❌ failed to fetch logs from both endpoints:\n• primary: %v\n• fallback: %v", err, err)
	}

	return fallbackResp, nil
}

func fetchWorkflowSteps(workflowName, token string) (string, []string, map[string]string, error) {
	url := fmt.Sprintf("http://localhost:2746/api/v1/workflows/cas/%s", workflowName)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Metadata struct {
			Uid string `json:"uid"`
		} `json:"metadata"`
		Status struct {
			Nodes map[string]struct {
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				Phase       string `json:"phase"`
			} `json:"nodes"`
		} `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", nil, nil, err
	}

	var names []string
	idMap := make(map[string]string)
	for id, node := range result.Status.Nodes {
		if node.Type == "Pod" {
			names = append(names, node.DisplayName)
			idMap[node.DisplayName] = id
		}
	}

	sort.Strings(names)
	return result.Metadata.Uid, names, idMap, nil
}
