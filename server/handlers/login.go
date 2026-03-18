package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"heartchat-server/services"
)

// LoginRequest is the payload for POST /api/login.
type LoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// LoginResponse is returned from POST /api/login.
// openID is used as a stable user identifier on all subsequent requests.
type LoginResponse struct {
	OpenID string `json:"openId"`
}

// Login handles POST /api/login.
// It exchanges the WeChat login code for the user's stable openid and returns
// it to the client so every subsequent request can include it as the user key.
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	openID, err := services.WeChatLogin(req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{OpenID: openID})
}
