package route

import (
	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/handler"
)

func SetupRoutes(router *gin.Engine, handler *handler.Handler, jwtSecret string) {
	SetupAuthRoutes(router, handler)
}
