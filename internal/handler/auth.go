package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/service/auth"
)

type AuthHandler struct {
	authService *auth.AuthService
	jwtSecret   string
}

func NewAuthHandler(authService *auth.AuthService, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtSecret:   jwtSecret,
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

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error{Message: "Invalid request data"})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error{Message: "Invalid request data"})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.Error{Message: err.Error()})
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusOK, model.Token{Value: token})
}
