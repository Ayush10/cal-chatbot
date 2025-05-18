package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

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

		// Cal.com email verification and events
		cal := api.Group("/cal")
		{
			cal.POST("/request-verification-code", h.HandleRequestVerificationCode)
			cal.POST("/verify-email-code", h.HandleVerifyEmailCode)
			cal.GET("/scheduled-events", h.HandleGetScheduledEvents)
		}
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

// HandleRequestVerificationCode proxies a request to Cal.com to send a verification code to the user's email
func (h *Handler) HandleRequestVerificationCode(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid email"})
		return
	}
	// Call Cal.com API
	calRes, err := proxyCalComRequest("POST", "/v2/verified-resources/emails/verification-code/request", map[string]string{"email": req.Email})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, calRes)
}

// HandleVerifyEmailCode proxies a request to Cal.com to verify the code for the user's email
func (h *Handler) HandleVerifyEmailCode(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" || req.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid email/code"})
		return
	}
	// Call Cal.com API
	calRes, err := proxyCalComRequest("POST", "/v2/verified-resources/emails/verification-code/verify", map[string]string{"email": req.Email, "code": req.Code})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	// If verification is successful, set a secure cookie with the verified email
	if status, ok := calRes["status"].(string); ok && status == "success" {
		c.SetCookie("verified_email", req.Email, 3600, "/", "", true, true)
	}
	c.JSON(http.StatusOK, calRes)
}

// HandleGetScheduledEvents proxies a request to Cal.com to get scheduled events for the verified email
func (h *Handler) HandleGetScheduledEvents(c *gin.Context) {
	email, err := c.Cookie("verified_email")
	if err != nil || email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not verified or missing email in session"})
		return
	}
	// Call Cal.com API for events (replace with correct endpoint as needed)
	calRes, err := proxyCalComRequest("GET", "/v2/bookings?email="+email, nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, calRes)
}

// proxyCalComRequest is a helper to call Cal.com API with your API key
func proxyCalComRequest(method, path string, body interface{}) (map[string]interface{}, error) {
	apiKey := os.Getenv("CALCOM_API_KEY")
	if apiKey == "" {
		return nil, gin.Error{Err: io.EOF, Type: gin.ErrorTypePrivate, Meta: "CALCOM_API_KEY not set"}
	}
	url := "https://api.cal.com" + path
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(b)
	}
	request, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+apiKey)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
