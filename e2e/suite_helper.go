package e2e

import (
	"os"
	// "testing"

	. "github.com/onsi/ginkgo/v2"
	// . "github.com/onsi/gomega"
)

// This file helps set up the test environment and can be expanded with 
// before/after suite hooks, environment preparation, etc.

var (
	// Add any global variables or test configurations here
	testServerURL string
)

var _ = BeforeSuite(func() {
	// Setup logic before running the entire test suite
	// For example, starting your server, setting up test databases, etc.
	testServerURL = os.Getenv("TEST_SERVER_URL")
	if testServerURL == "" {
		testServerURL = "http://localhost:8080"
	}

	// You can add more setup logic here
	By("Ensuring test environment is ready")
	// Add any necessary checks or preparatory steps
})

var _ = AfterSuite(func() {
	// Cleanup logic after running the entire test suite
	// For example, stopping servers, cleaning up resources
	By("Cleaning up test environment")
	// Add any necessary cleanup steps
})