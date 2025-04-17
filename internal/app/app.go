package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/config"
	"github.com/kirillidk/pvz-service/internal/handler"
	"github.com/kirillidk/pvz-service/internal/middleware"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/pkg/database"
)

type App struct {
	config  *config.Config
	router  *gin.Engine
	handler *handler.Handler
}

func NewApp(cfg *config.Config) *App {
	return &App{
		config:  cfg,
		router:  gin.Default(),
		handler: handler.NewHandler(cfg),
	}
}

func (a *App) Run() error {
	_, err := database.NewPostgresDB(&a.config.Database)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	a.setupRoutes()

	serverAddr := fmt.Sprintf(":%s", a.config.Server.Port)
	log.Printf("Starting server on %s", serverAddr)

	return a.router.Run(serverAddr)
}

func (a *App) setupRoutes() {

	a.router.POST("/dummyLogin", a.handler.AuthHandler.DummyLogin)

	a.router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "lol")
	})

	authenticated := a.router.Group("/", middleware.AuthMiddleware(a.config.JWT.JWTSecret))

	employees := authenticated.Group("/", middleware.RoleMiddleware(model.EmployeeRole))
	employees.GET("/emp", func(c *gin.Context) {
		c.String(http.StatusOK, "emp")
	})

	moderators := authenticated.Group("/", middleware.RoleMiddleware(model.ModeratorRole))
	moderators.GET("/mod", func(c *gin.Context) {
		c.String(http.StatusOK, "mod")
	})
}
