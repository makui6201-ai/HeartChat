package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"heartchat-server/services"
)

// HistoryResponse is returned from GET /api/history.
type HistoryResponse struct {
	Role    string             `json:"role"`
	Love    int                `json:"love"`
	History []services.Message `json:"history"`
}

// GetHistory handles GET /api/history?openId=<id>.
// Returns the stored role, love score, and conversation history for the user.
func GetHistory(c *gin.Context) {
	openID := c.Query("openId")
	if openID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "openId is required"})
		return
	}

	state, err := services.Load(openID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, HistoryResponse{
		Role:    state.Role,
		Love:    state.Love,
		History: state.History,
	})
}

// DeleteHistory handles DELETE /api/history.
// Clears the conversation history for the user while keeping role and love score.
type DeleteHistoryRequest struct {
	OpenID string `json:"openId" binding:"required"`
}

func DeleteHistory(c *gin.Context) {
	var req DeleteHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ClearHistory(req.OpenID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
