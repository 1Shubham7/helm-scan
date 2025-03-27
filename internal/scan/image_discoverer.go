package scan

import (
	// "context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/1shubham7/helm-scan/internal/models"
	// "github.com/docker/docker/client"
)

func DiscoverImages(extractedChartPath string) ([]models.ImageInfo, error) {

	patterns := []*regexp.Regexp{
		// Matches: image: "nginx:1.21"  OR  image: postgres:13
		// Explanation:
		// - `image:\s*` → Matches "image:" followed by any spaces.
		// - `["']?` → Optionally matches a single or double quote around the image name.
		// - `([^"'\s]+)` → Captures the actual image name (anything except spaces or quotes).
		// - `["']?` → Optionally matches the closing quote.
		// regexp.MustCompile(`image:\s*["']?([^"'\s]+)["']?`),

		// Matches: repository: "myrepo/custom-app:v2"  OR  repository: myrepo/custom-app:v2
		// Explanation:
		// - `repository:\s*` → Matches "repository:" followed by any spaces.
		// - `["']?` → Optionally matches a single or double quote around the repository name.
		// - `([^"'\s]+)` → Captures the repository name.
		// - `["']?` → Optionally matches the closing quote.
		regexp.MustCompile(`(?:image|repository):\s*["']?([a-zA-Z0-9.-]+(?:/[a-zA-Z0-9.-]+)*(?::[a-zA-Z0-9.-]+)?)["']?`),
	}

	// store images
	images := make(map[string]models.ImageInfo)

	err := filepath.Walk(extractedChartPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip dirs and only processs yaml files
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		for _, pattern := range patterns {
			matches := pattern.FindAllStringSubmatch(string(content), -1)
			for _, match := range matches {
				if len(match) > 1 {
					imageName := match[1]

					if isValidImageReference(imageName) {

						// Normalize image reference
						normalizedImage := normalizeImage(imageName)

						// Store unique image
						images[normalizedImage.Name] = normalizedImage
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk chart directory: %w", err)
	}

	// converting map to slice
	respImages := make([]models.ImageInfo, 0, len(images))
	for _, i := range images {
		respImages = append(respImages, i)
	}

	return respImages, nil
}

// normalizeImage standardizes image references
func normalizeImage(rawImage string) models.ImageInfo {
	// Default tag if not specified
	defaultTag := "latest"

	// Split image into parts
	parts := strings.Split(rawImage, ":")
	name := parts[0]
	tag := defaultTag

	// Set tag if provided
	if len(parts) > 1 {
		tag = parts[1]
	}

	// Determine repository
	repository := ""
	if strings.Contains(name, "/") {
		// If image name contains '/', assume it's a full path
		repository = strings.Join(strings.Split(name, "/")[:len(strings.Split(name, "/"))-1], "/")
	}

	// // rawImage is image + tag
	size, layers, _ := getSizeAndLayers(rawImage)

	return models.ImageInfo{
		Name:       name,
		Repository: repository,
		Tag:        tag,
		Size:       size,
		Layers:     layers,
	}
}

func getSizeAndLayers(imageWithTag string) (size int64, layers int, err error) {
	size = 0
	layers = 0

	// Run docker inspect to get detailed image information
	cmd := exec.Command("docker", "inspect",
		"--format={{.Size}}\t{{len .RootFS.Layers}}",
		imageWithTag)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("galati")
	}

	// Trim any whitespace
	outputStr := strings.TrimSpace(string(output))

	// Split the output into size and layers
	parts := strings.Split(outputStr, "\t")
	if len(parts) < 2 {
		return 0, 0, fmt.Errorf("unexpected output format for image %s", imageWithTag)
	}

	// Parse size
	sizeVal, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse size for image %s: %w", imageWithTag, err)
	}
	size = sizeVal

	// Parse layers
	layersVal, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse layers for image %s: %w", imageWithTag, err)
	}
	layers = layersVal

	return size, layers, nil
}

// func getSizeAndLayers(imageWithTag string) (size int64, layers int, err error) {
// 	ctx := context.Background()
// 	// my Docker daemon supports 1.47 at most, and client was of latest 1.48 version.
// 	// cli, err := client.NewClientWithOpts(client.FromEnv)

// 	// this give client highest API version that both the client and daemon support.
// 	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 	if err != nil {
// 		return 0, 0, fmt.Errorf("failed to create Docker client: %w", err)
// 	}
// 	defer cli.Close()

// 	inspect, err := cli.ImageInspect(ctx, imageWithTag)
// 	if err != nil {
// 		fmt.Println("one", imageWithTag)
// 		return 0, 0, fmt.Errorf("failed to inspect image %s: %w", imageWithTag, err)
// 	}

// 	size = inspect.Size
// 	layers = len(inspect.RootFS.Layers)

// 	fmt.Println(inspect.RootFS.Layers)

// 	return size, layers, nil
// }

// isValidImageReference checks if the image reference is potentially valid
func isValidImageReference(image string) bool {
	// Checks to filter out invalid references
	invalidPatterns := []string{
		`^\{\{.*\}\}$`,                       // Helm template variables
		`^["']+$`,                            // Just quotes
		`^[{}]+$`,                            // Just braces
		`^\s*$`,                              // Empty or whitespace
		`^[/:.]+$`,                           // Just delimiters
		`^(your-image|placeholder|example)$`, // Common placeholder names
		`^(oci|registry)$`,                   // Generic registry names
	}

	for _, pattern := range invalidPatterns {
		match, _ := regexp.MatchString(pattern, image)
		if match {
			return false
		}
	}

	// Additional check: must contain alphanumeric characters
	return regexp.MustCompile(`[a-zA-Z0-9]`).MatchString(image)
}
