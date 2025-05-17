package api

import (
	"net/http"

	"log"

	"github.com/gin-gonic/gin"
)

// HandleHealth handles health check requests
func (h *Handler) HandleHealth(c *gin.Context) {
	log.Printf("[INFO] Health check requested from %s", c.ClientIP())
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}
