package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ImageDetailsHandler(c *gin.Context) {
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