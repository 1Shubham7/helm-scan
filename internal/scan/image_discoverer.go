package scan

// import (
// 	"regexp"

// 	"github.com/1shubham7/helm-scan/internal/models"
// )

// func DiscoverImages(chartPath string) ([]models.ImageInfo, error) {

// 	patterns := []*regexp.Regexp{
// 		// Matches: image: "nginx:1.21"  OR  image: postgres:13  
// 	// Explanation:
// 	// - `image:\s*` → Matches "image:" followed by any spaces.
// 	// - `["']?` → Optionally matches a single or double quote around the image name.
// 	// - `([^"'\s]+)` → Captures the actual image name (anything except spaces or quotes).
// 	// - `["']?` → Optionally matches the closing quote.
// 	regexp.MustCompile(`image:\s*["']?([^"'\s]+)["']?`),

// 	// Matches: repository: "myrepo/custom-app:v2"  OR  repository: myrepo/custom-app:v2
// 	// Explanation:
// 	// - `repository:\s*` → Matches "repository:" followed by any spaces.
// 	// - `["']?` → Optionally matches a single or double quote around the repository name.
// 	// - `([^"'\s]+)` → Captures the repository name.
// 	// - `["']?` → Optionally matches the closing quote.
// 	regexp.MustCompile(`repository:\s*["']?([^"'\s]+)["']?`),
// 	}
// }