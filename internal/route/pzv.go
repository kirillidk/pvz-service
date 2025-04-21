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

		pvzGroup.GET("", middleware.RoleMiddleware(model.EmployeeRole, model.ModeratorRole), handler.PVZHandler.GetPVZList)

		pvzGroup.POST("", middleware.RoleMiddleware(model.ModeratorRole), handler.PVZHandler.CreatePVZ)
		pvzGroup.POST("/:pvzId/delete_last_product", middleware.RoleMiddleware(model.EmployeeRole), handler.ProductHandler.DeleteLastProduct)
		pvzGroup.POST("/:pvzId/close_last_reception", middleware.RoleMiddleware(model.EmployeeRole), handler.ReceptionHandler.CloseLastReception)
	}
}
