package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"heartchat-server/services"
)

// ChatRequest is the payload for POST /api/chat.
type ChatRequest struct {
	OpenID  string             `json:"openId"`
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

	// If an openId is provided, load persisted state.  Server-side history is
	// used as a fallback when the client sends no history (e.g. on a fresh
	// page load); if the client already sent history it is used as-is.
	history := req.History
	if req.OpenID != "" {
		state, err := services.Load(req.OpenID)
		if err != nil {
			log.Printf("load user state error: %v", err)
		} else {
			// Update persisted role/love to reflect the current request.
			state.Role = req.Role
			state.Love = req.Love
			if err := services.Save(req.OpenID, state); err != nil {
				log.Printf("save user state error: %v", err)
			}
			// Use server-side history when the client didn't send any.
			if len(history) == 0 && len(state.History) > 0 {
				history = state.History
			}
		}
	}

	// Keep at most the last 10 history messages to stay within context limits.
	if len(history) > 10 {
		history = history[len(history)-10:]
	}

	reply, err := services.GetAIReply(req.Role, req.Love, req.Message, history)
	if err != nil {
		log.Printf("GetAIReply error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI服务暂时不可用，请稍后再试"})
		return
	}

	// Persist the new turn to server-side history.
	if req.OpenID != "" {
		if err := services.AppendHistory(req.OpenID, req.Message, reply); err != nil {
			log.Printf("append history error: %v", err)
		}
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
