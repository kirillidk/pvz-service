package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/config"
	grpcserver "github.com/kirillidk/pvz-service/internal/grpc"
	"github.com/kirillidk/pvz-service/internal/handler"
	"github.com/kirillidk/pvz-service/internal/repository"
	"github.com/kirillidk/pvz-service/internal/route"
	"github.com/kirillidk/pvz-service/internal/service"
	grpcservice "github.com/kirillidk/pvz-service/internal/service/grpc"
	"github.com/kirillidk/pvz-service/pkg/database"
)

type App struct {
	Config     *config.Config
	Router     *gin.Engine
	Database   *sql.DB
	Repository *repository.Repository
	Service    *service.Service
	Handler    *handler.Handler
	GRPCServer *grpcserver.Server
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

	grpcPVZService := grpcservice.NewPVZService(repo.PVZRepository)
	grpcSrv := grpcserver.NewServer(cfg, grpcPVZService)

	return &App{
		Config:     cfg,
		Database:   db,
		Router:     rtr,
		Repository: repo,
		Service:    serv,
		Handler:    handl,
		GRPCServer: grpcSrv,
	}, nil
}

func (a *App) Run() error {
	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.GRPCServer.Start(); err != nil {
			log.Printf("gRPC server error: %v", err)
			errCh <- err
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		serverAddr := fmt.Sprintf(":%s", a.Config.Server.Port)
		log.Printf("Starting HTTP server on %s", serverAddr)

		if err := a.Router.Run(serverAddr); err != nil {
			log.Printf("HTTP server error: %v", err)
			errCh <- err
			cancel()
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
