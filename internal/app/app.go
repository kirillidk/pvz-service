package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/config"
	"github.com/kirillidk/pvz-service/internal/handler"
	"github.com/kirillidk/pvz-service/internal/middleware"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/repository"
	"github.com/kirillidk/pvz-service/internal/service"
	"github.com/kirillidk/pvz-service/pkg/database"
)

type App struct {
	config     *config.Config
	router     *gin.Engine
	database   *sql.DB
	repository *repository.Repository
	service    *service.Service
	handler    *handler.Handler
}

func NewApp(cfg *config.Config) (*App, error) {
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	rtr := gin.Default()
	repo := repository.NewRepository(db)
	serv := service.NewService(repo, cfg.JWT.JWTSecret)
	handl := handler.NewHandler(serv, cfg)

	return &App{
		config:     cfg,
		database:   db,
		router:     rtr,
		repository: repo,
		service:    serv,
		handler:    handl,
	}, nil
}

func (a *App) Run() error {
	a.setupRoutes()

	serverAddr := fmt.Sprintf(":%s", a.config.Server.Port)
	log.Printf("Starting server on %s", serverAddr)

	return a.router.Run(serverAddr)
}

func (a *App) setupRoutes() {

	a.router.POST("/dummyLogin", a.handler.AuthHandler.DummyLogin)
	a.router.POST("/register", a.handler.AuthHandler.Register)
	a.router.POST("/login", a.handler.AuthHandler.Login)

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
