package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// ScanRequest represents the structure of the scan request
type ScanRequest struct {
	ChartURL string `json:"chartURL"`
}

// ScanResponse represents the structure of the scan response
type ScanResponse struct {
	Images []struct {
		// Define the structure of your Images slice if needed
		// For now, we'll keep it generic
	} `json:"images"`
}

var _ = Describe("Chart Scan API", func() {
	It("should successfully scan a chart and return 200 status", func() {
		// Prepare the request body
		requestBody := ScanRequest{
			ChartURL: "oci://registry-1.docker.io/bitnamicharts/nats",
		}

		// Convert request body to JSON
		jsonBody, err := json.Marshal(requestBody)
		Expect(err).NotTo(HaveOccurred())

		// Send POST request to the scan endpoint
		resp, err := http.Post("http://localhost:8080/scan", "application/json", bytes.NewBuffer(jsonBody))
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()

		// Check status code
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		// Parse response body
		var scanResponse ScanResponse
		err = json.NewDecoder(resp.Body).Decode(&scanResponse)
		Expect(err).NotTo(HaveOccurred())

		// Optional: You might want to add more specific assertions about the response
		// For example, checking if Images slice is not empty
		Expect(scanResponse.Images).NotTo(BeEmpty())
	})
})

// This function is required by Ginkgo to run the tests
func TestChartScanAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chart Scan API Suite")
}