package app

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/config"
	"github.com/kirillidk/pvz-service/pkg/database"
)

type App struct {
	config *config.Config
	router *gin.Engine
}

func NewApp(cfg *config.Config) *App {
	return &App{
		config: cfg,
		router: gin.Default(),
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
	a.router.GET("/hello", func(c *gin.Context) {
		c.String(200, "lol")
	})
}
