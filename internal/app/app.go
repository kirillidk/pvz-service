package app

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/config"
	"github.com/kirillidk/pvz-service/internal/handler"
	"github.com/kirillidk/pvz-service/internal/repository"
	"github.com/kirillidk/pvz-service/internal/route"
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

	repo := repository.NewRepository(db)
	serv := service.NewService(repo, cfg.JWT.JWTSecret)
	handl := handler.NewHandler(serv, cfg)

	rtr := gin.Default()
	route.SetupRoutes(rtr, handl, cfg.JWT.JWTSecret)

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
	serverAddr := fmt.Sprintf(":%s", a.config.Server.Port)
	log.Printf("Starting server on %s", serverAddr)

	return a.router.Run(serverAddr)
}
