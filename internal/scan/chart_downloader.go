package scan

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func Download(chartURL string) (string, error) {
	err := validateChartURL(chartURL)
	if err != nil {
		return "", fmt.Errorf("failed to create pull the chart: %w", err)
	}

	// create temp dir
	tempDir, err := os.MkdirTemp("", "helm-chart-*")
	// "" will create dir in default temp directory of the system
	// Linux/macOS: /tmp/
	// Windows: C:\Users\Username\AppData\Local\Temp\
	// tempDir = /tmp/helm-chart-abc123
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	return pullChart(chartURL, tempDir)
}

func pullChart(chartURL, tempDir string) (string, error) {

	cmd := exec.Command("helm", "pull",
		chartURL,
		"--untar",
		"--destination", tempDir,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to pull Helm chart: %v", err)
	}

	return tempDir, nil
}

func validateChartURL(chartURL string) error {
	// Check if URL is empty
	if chartURL == "" {
		return fmt.Errorf("chart URL cannot be empty")
	}

	// Check for supported URL prefixes
	supportedPrefixes := []string{
		"oci://",
		"https://",
		"http://",
		"file://",
	}

	hasValidPrefix := false
	for _, prefix := range supportedPrefixes {
		if strings.HasPrefix(chartURL, prefix) {
			hasValidPrefix = true
			break
		}
	}

	if !hasValidPrefix {
		return fmt.Errorf("invalid chart URL prefix. Supported prefixes: %v", supportedPrefixes)
	}

	// For OCI URLs, do additional validation
	if strings.HasPrefix(chartURL, "oci://") {
		// Validate OCI registry URL structure
		registryParts := strings.Split(strings.TrimPrefix(chartURL, "oci://"), "/")
		if len(registryParts) < 2 {
			return fmt.Errorf("invalid OCI chart URL format. Expected: oci://registry/repository/chart")
		}
	}

	// Optional: Parse the URL to do more detailed validation
	parsedURL, err := url.Parse(chartURL)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	// Additional checks can be added based on your specific requirements
	// For example, checking against a whitelist of known registries
	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL must have a valid scheme")
	}

	return nil
}
