package scan

import (
	"fmt"
	"os"
	"os/exec"
)

func Download(chartURL string) (string, error) {
	// validate url

	// create temp dir
	tempDir, err := os.MkdirTemp("", "helm-chart-*")
	// "" will create dir in default temp directory of the system
	// Linux/macOS: /tmp/
	// Windows: C:\Users\Username\AppData\Local\Temp\
	// tempDir = /tmp/helm-chart-abc123
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	return downloadChart(chartURL, tempDir)
}

func downloadChart(chartURL, tempDir string) (string, error) {

	cmd := exec.Command("helm", "pull",
		chartURL,                 // Chart reference (e.g., oci://registry-1.docker.io/bitnamicharts/airflow)
		"--untar",                // Automatically untar the chart
		"--destination", tempDir, // Extract to current directory
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to pull Helm chart: %v", err)
	}

	return tempDir, nil
}
