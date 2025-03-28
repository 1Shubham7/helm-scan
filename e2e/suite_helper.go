package e2e

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
)

var (
	testServerURL string
)

var _ = BeforeSuite(func() {
	testServerURL = os.Getenv("TEST_SERVER_URL")
	if testServerURL == "" {
		testServerURL = "http://localhost:8080"
	}

	By("Ensuring test environment is ready")
})

var _ = AfterSuite(func() {
	By("Cleaning up test environment")
})
