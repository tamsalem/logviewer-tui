package argo

import (
	"encoding/base64"
	"fmt"
	"os/exec"
)

func getArgoToken() (string, error) {
	// Try argo-server token first
	token, err := getTokenFromSecret("argo-server.service-account-token")
	if err == nil && token != "" {
		return "Bearer " + token, nil
	}

	// Fallback to argo-controller for Local Env
	token, err = getTokenFromSecret("argo-controller.service-account-token")
	if err == nil && token != "" {
		return "Bearer " + token, nil
	}

	return "", fmt.Errorf("failed to retrieve Argo token from argo-server and argo-controller")
}

func getTokenFromSecret(secretName string) (string, error) {
	cmd := exec.Command("kubectl", "-n", "cas", "get", "secret", secretName, "-o", "jsonpath={.data.token}")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	decoded, err := base64.StdEncoding.DecodeString(string(output))
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
