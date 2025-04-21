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
	Config     *config.Config
	Router     *gin.Engine
	Database   *sql.DB
	Repository *repository.Repository
	Service    *service.Service
	Handler    *handler.Handler
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
		Config:     cfg,
		Database:   db,
		Router:     rtr,
		Repository: repo,
		Service:    serv,
		Handler:    handl,
	}, nil
}

func (a *App) Run() error {
	serverAddr := fmt.Sprintf(":%s", a.Config.Server.Port)
	log.Printf("Starting server on %s", serverAddr)

	return a.Router.Run(serverAddr)
}
