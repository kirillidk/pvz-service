package route

import (
	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/handler"
	"github.com/kirillidk/pvz-service/internal/middleware"
	"github.com/kirillidk/pvz-service/internal/model"
)

func SetupReceptionRoutes(router *gin.Engine, handler *handler.Handler, jwtSecret string) {
	receptionGroup := router.Group("/receptions")
	{
		receptionGroup.Use(middleware.AuthMiddleware(jwtSecret))

		receptionGroup.POST("", middleware.RoleMiddleware(model.EmployeeRole), handler.ReceptionHandler.CreateReception)
	}
}
