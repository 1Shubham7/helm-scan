package api

import (
	"net/http"

	"github.com/1shubham7/helm-scan/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type scanRequest struct {
	ChartURL string `json:"chartURL"`
}

type ScanResponse struct {
	Images []models.ImageInfo `json:"images"`
}

var validate = validator.New()

func ChartScanHandler(c *gin.Context){
	var req scanRequest

	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request is invalid",
			"details": err.Error(),
		})
		return
	}

	err = validate.Struct(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request is not valid",
			"details": err.Error(),
		})
	}

	// validate the chart link
	if req.ChartURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request is invalid",
			"details": "Chart URL is required",
		})
		return
	}

	images, err := ScanChart(scanRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to scan",
			"details": err.Error(),
		})
	}

	c.JSON(http.StatusOK, ScanResponse{
		Images: images,
	})
}