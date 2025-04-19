package route

import (
	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/handler"
)

func SetupAuthRoutes(router *gin.Engine, handler *handler.Handler) {
	router.POST("/dummyLogin", handler.AuthHandler.DummyLogin)
	router.POST("/register", handler.AuthHandler.Register)
	router.POST("/login", handler.AuthHandler.Login)
}
