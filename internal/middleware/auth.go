package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/model"
	service "github.com/kirillidk/pvz-service/internal/service/auth"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, model.Error{Message: "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, model.Error{Message: "Authorization header must be in format: Bearer {token}"})
			c.Abort()
			return
		}

		claims, err := service.ValidateToken(parts[1], jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.Error{Message: "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("userRole", claims.Role)

		c.Next()
	}
}

func RoleMiddleware(requiredRoles ...model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.Error{Message: "User not authenticated"})
			c.Abort()
			return
		}

		for _, role := range requiredRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, model.Error{Message: "Operation not permitted for this user role"})
		c.Abort()
	}
}
