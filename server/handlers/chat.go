package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"heartchat-server/services"
)

// ChatRequest is the payload for POST /api/chat.
type ChatRequest struct {
	Role    string             `json:"role"    binding:"required,oneof=boyfriend girlfriend"`
	Love    int                `json:"love"    binding:"min=0"`
	Message string             `json:"message" binding:"required,max=500"`
	History []services.Message `json:"history"`
}

// ChatResponse is returned from POST /api/chat.
type ChatResponse struct {
	Reply string `json:"reply"`
}

// GreetingRequest is the payload for POST /api/greeting.
type GreetingRequest struct {
	Role string `json:"role" binding:"required,oneof=boyfriend girlfriend"`
	Love int    `json:"love" binding:"min=0"`
}

// GreetingResponse is returned from POST /api/greeting.
type GreetingResponse struct {
	Greeting string `json:"greeting"`
}

// Chat handles POST /api/chat.
func Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Keep at most the last 10 history messages to stay within context limits.
	if len(req.History) > 10 {
		req.History = req.History[len(req.History)-10:]
	}

	reply, err := services.GetAIReply(req.Role, req.Love, req.Message, req.History)
	if err != nil {
		log.Printf("GetAIReply error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI服务暂时不可用，请稍后再试"})
		return
	}

	c.JSON(http.StatusOK, ChatResponse{Reply: reply})
}

// Greeting handles POST /api/greeting.
func Greeting(c *gin.Context) {
	var req GreetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GreetingResponse{Greeting: services.GetGreeting(req.Role, req.Love)})
}
