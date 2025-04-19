package route

import (
	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/handler"
	"github.com/kirillidk/pvz-service/internal/middleware"
	"github.com/kirillidk/pvz-service/internal/model"
)

func SetupPVZRoutes(router *gin.Engine, handler *handler.Handler, jwtSecret string) {
	pvzGroup := router.Group("/pvz")
	{
		pvzGroup.Use(middleware.AuthMiddleware(jwtSecret))

		pvzGroup.POST("", middleware.RoleMiddleware(model.ModeratorRole), handler.PVZHandler.CreatePVZ)
	}
}
