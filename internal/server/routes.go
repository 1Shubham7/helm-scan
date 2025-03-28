package server

import (
	"net/http"

	api "github.com/1shubham7/helm-scan/api"
	"github.com/gin-gonic/gin"

	"github.com/1shubham7/helm-scan/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	router := gin.Default()

	router.Use(middleware.PerClientTokenBucket())

	router.GET("/", s.HelloWorldHandler)
	router.POST("/scan", api.ChartScanHandler)

	return router
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}
