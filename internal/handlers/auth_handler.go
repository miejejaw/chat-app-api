package handlers

import (
	"chat-app-api/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.BindJSON(&loginRequest); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	response, err := h.authService.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	ctx.JSON(200, gin.H{
		"data":    response,
		"message": "Login successful",
	})
}

func (h *AuthHandler) RenewAccessToken(ctx *gin.Context) {
	var renewRequest struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := ctx.BindJSON(&renewRequest); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	newAccessToken, err := h.authService.RenewAccessToken(renewRequest.RefreshToken)
	if err != nil {
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	ctx.JSON(200, gin.H{
		"access_token": newAccessToken,
	})
}
