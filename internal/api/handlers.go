package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/cal-chatbot/internal/chatbot"
)

// Handler contains all API handlers
type Handler struct {
	chatbot *chatbot.Chatbot
}

// NewHandler creates a new API handler
func NewHandler(bot *chatbot.Chatbot) *Handler {
	return &Handler{
		chatbot: bot,
	}
}

// SetupRoutes sets up the API routes
func (h *Handler) SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/chat", h.HandleChat)
		api.GET("/health", h.HandleHealth)
		api.GET("/endpoints", h.HandleEndpoints)
		api.GET("/history/:conversation_id", h.HandleLoadHistory)
		api.GET("/history/search", h.HandleSearchHistory)
	}

	// Serve static files for the web interface
	r.Static("/web", "./web")
	r.StaticFile("/", "./web/index.html")
}

// HandleEndpoints returns a list of all available API endpoints
func (h *Handler) HandleEndpoints(c *gin.Context) {
	log.Printf("[INFO] Endpoints listed for %s", c.ClientIP())
	c.JSON(200, gin.H{
		"endpoints": []string{
			"POST /api/chat",
			"GET /api/health",
			"GET /api/endpoints",
		},
	})
}

// HandleLoadHistory loads a conversation's history
func (h *Handler) HandleLoadHistory(c *gin.Context) {
	conversationID := c.Param("conversation_id")
	history, err := chatbot.LoadHistory(conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"history": history})
}

// HandleSearchHistory searches all conversations for a term
func (h *Handler) HandleSearchHistory(c *gin.Context) {
	term := c.Query("q")
	if term == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing search term."})
		return
	}
	matches, err := chatbot.SearchHistory(term)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"matches": matches})
}
