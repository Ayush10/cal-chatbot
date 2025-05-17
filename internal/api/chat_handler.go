package api

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/cal-chatbot/internal/chatbot"
	"github.com/yourusername/cal-chatbot/internal/models"
)

// logError logs errors with context and request ID
func logError(context, requestID string, err error) {
	if requestID != "" {
		log.Printf("[ERROR] [%s] %s: %v", requestID, context, err)
	} else {
		log.Printf("[ERROR] %s: %v", context, err)
	}
}

// HandleChat handles chat messages
func (h *Handler) HandleChat(c *gin.Context) {
	// Generate or extract a conversation ID for history
	conversationID := c.GetHeader("X-Conversation-Id")
	if conversationID == "" {
		conversationID = uuid.New().String()
	}

	// Debug: print raw request body and headers
	log.Printf("[DEBUG] [%s] Content-Type: %s", conversationID, c.GetHeader("Content-Type"))
	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	log.Printf("[DEBUG] [%s] Raw request body: %s", conversationID, string(bodyBytes))
	// Re-parse the body for binding
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	var req models.ChatRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logError("JSON binding error", conversationID, err)
		log.Printf("[DEBUG] [%s] Failed to bind JSON. Raw body: %s", conversationID, string(bodyBytes))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Your request was not understood. Please check your input and try again.",
		})
		return
	}
	log.Printf("[DEBUG] [%s] Parsed ChatRequest: %+v", conversationID, req)
	log.Printf("[DEBUG] [%s] req.Messages type: %T, len: %d", conversationID, req.Messages, len(req.Messages))

	if len(req.Messages) == 0 {
		logError("Empty messages array", conversationID, nil)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please provide at least one message in your request.",
		})
		return
	}

	// Save all user messages to history
	for _, m := range req.Messages {
		if m.Role == "user" && m.Content != "" {
			_ = chatbot.SaveMessage(conversationID, m.Role, m.Content)
		}
	}

	response, err := h.chatbot.ProcessMessage(c.Request.Context(), req.Messages)
	if err != nil {
		logError("Failed to process message", conversationID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Sorry, something went wrong while processing your message. Please try again later.",
		})
		return
	}

	// Save assistant response to history
	_ = chatbot.SaveMessage(conversationID, "assistant", response)

	c.Header("X-Conversation-Id", conversationID)
	c.JSON(http.StatusOK, models.ChatResponse{
		Message: response,
		// Optionally, add conversationID to the response for the frontend
		// ConversationID: conversationID,
	})
}
