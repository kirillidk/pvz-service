package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/service/auth"
)

type AuthHandler struct {
	jwtSecret string
}

func NewAuthHandler(jwtSecret string) *AuthHandler {
	return &AuthHandler{
		jwtSecret: jwtSecret,
	}
}

func (authHandler *AuthHandler) DummyLogin(c *gin.Context) {
	var req dto.DummyLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error{Message: "Invalid role. Must be 'employee' or 'moderator'"})
		return
	}

	token, err := auth.GenerateToken(req.Role, authHandler.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Message: "Failed to generate token"})
		return
	}

	c.Header("Authorization", "Bearer "+token)

	c.JSON(http.StatusOK, model.Token{Value: token})
}
