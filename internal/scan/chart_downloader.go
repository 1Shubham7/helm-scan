package scan

import (
	"fmt"
    // "log"
    "os"
    "os/exec"
    "path/filepath"
    "strings"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	// "helm.sh/helm/v3/pkg/downloader"
	// "helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/registry"
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

	if registry.IsOCI(chartURL) {
		return downloadOCIChart(chartURL, tempDir)
	}

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
    // Construct the Helm pull command
    cmd := exec.Command("helm", "pull", chartURL, "--destination", tempDir)
    
    // Run the command
    _, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("failed to download charttt: %w", err)
    }

    // Get the filename of the pulled chart
    files, err := os.ReadDir(tempDir)
    if err != nil {
        return "", fmt.Errorf("failed to read temporary directory: %w", err)
    }

    // Find the first .tgz file
    for _, file := range files {
        if strings.HasSuffix(file.Name(), ".tgz") {
            return filepath.Join(tempDir, file.Name()), nil
        }
    }

    return "", fmt.Errorf("no chart file found after download")
}