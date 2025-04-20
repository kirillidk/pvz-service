package route

import (
	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/handler"
	"github.com/kirillidk/pvz-service/internal/middleware"
	"github.com/kirillidk/pvz-service/internal/model"
)

func SetupProductRoutes(router *gin.Engine, handler *handler.Handler, jwtSecret string) {
	productGroup := router.Group("/products")
	{
		productGroup.Use(middleware.AuthMiddleware(jwtSecret))

		productGroup.POST("", middleware.RoleMiddleware(model.EmployeeRole), handler.ProductHandler.CreateProduct)
	}
}
