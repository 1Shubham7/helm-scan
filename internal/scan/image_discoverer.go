package scan

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/1shubham7/helm-scan/internal/models"
)

func DiscoverImages(extractedChartPath string) ([]models.ImageInfo, error) {

	patterns := []*regexp.Regexp{
		// Matches: image: "nginx:1.21"  OR  image: postgres:13  
	// Explanation:
	// - `image:\s*` → Matches "image:" followed by any spaces.
	// - `["']?` → Optionally matches a single or double quote around the image name.
	// - `([^"'\s]+)` → Captures the actual image name (anything except spaces or quotes).
	// - `["']?` → Optionally matches the closing quote.
	regexp.MustCompile(`image:\s*["']?([^"'\s]+)["']?`),

	// Matches: repository: "myrepo/custom-app:v2"  OR  repository: myrepo/custom-app:v2
	// Explanation:
	// - `repository:\s*` → Matches "repository:" followed by any spaces.
	// - `["']?` → Optionally matches a single or double quote around the repository name.
	// - `([^"'\s]+)` → Captures the repository name.
	// - `["']?` → Optionally matches the closing quote.
	regexp.MustCompile(`repository:\s*["']?([^"'\s]+)["']?`),
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

					// Normalize image reference
					normalizedImage := normalizeImage(imageName)
					
					// Store unique image
					images[normalizedImage.Name] = normalizedImage
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
	for  _, i := range images {
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
	repository := "docker.io/library"
	if strings.Contains(name, "/") {
		// If image name contains '/', assume it's a full path
		repository = strings.Join(strings.Split(name, "/")[:len(strings.Split(name, "/"))-1], "/")
	}

	return models.ImageInfo{
		Name:       name,
		Repository: repository,
		Tag:        tag,
	}
}