package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type ScanRequest struct {
	ChartURL string `json:"chartURL"`
}

type ScanResponse struct {
	Images []struct{} `json:"images"`
}

var _ = Describe("Chart Scan API", func() {
	It("should successfully scan a chart and return 200 status", func() {
		requestBody := ScanRequest{
			ChartURL: "oci://registry-1.docker.io/bitnamicharts/nats",
		}

		jsonBody, err := json.Marshal(requestBody)
		Expect(err).NotTo(HaveOccurred())

		// POST requires a io.Reader as the request body. jsonBody is a []byte, NewBuffer returns a new buffer which implements io.Reader
		jsonReader := bytes.NewBuffer(jsonBody)
		resp, err := http.Post("http://localhost:8080/scan", "application/json", jsonReader)
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()

		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		var scanResponse ScanResponse
		err = json.NewDecoder(resp.Body).Decode(&scanResponse)
		Expect(err).NotTo(HaveOccurred())

		Expect(scanResponse.Images).NotTo(BeEmpty())
	})
})

func TestChartScanAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chart Scan API Suite")
}
