package scan

import (
	"fmt"
    // "log"
    "os"
    "os/exec"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	// "helm.sh/helm/v3/pkg/downloader"
	// "helm.sh/helm/v3/pkg/getter"
	// "helm.sh/helm/v3/pkg/registry"
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

	// if registry.IsOCI(chartURL) {
		// return downloadOCIChart(chartURL, tempDir)
	// }

	return downloadHTTPChart(chartURL, tempDir)
}

func downloadOCIChart(chartURL, tempDir string) (string, error) {
	actionConfig := new(action.Configuration)
	settings := cli.New()

	if err := actionConfig.Init(settings.RESTClientGetter(), "", "", nil); err != nil {
		return "", fmt.Errorf("failed to initialize Helm action config: %w", err)
	}

	// Pull the OCI chart
	pullAction := action.NewPull()
	pullAction.Settings = settings
	pullAction.DestDir = tempDir
	pullAction.Untar = true
	pullAction.UntarDir = tempDir

	// Pull the chart
	chartPath, err := pullAction.Run(chartURL)
	if err != nil {
		return "", fmt.Errorf("failed to pull OCI chart: %w", err)
	}

	return chartPath, nil
}

func downloadHTTPChart(chartURL, tempDir string) (string, error){

	cmd := exec.Command("helm", "pull", 
		chartURL,       // Chart reference (e.g., oci://registry-1.docker.io/bitnamicharts/airflow)
		"--untar",      // Automatically untar the chart
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