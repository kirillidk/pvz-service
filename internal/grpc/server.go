package grpc

import (
	"fmt"
	"log"
	"net"

	pvz_v1 "github.com/kirillidk/pvz-service/api/proto/pvz/pvz_v1"
	"github.com/kirillidk/pvz-service/internal/config"
	grpcservice "github.com/kirillidk/pvz-service/internal/service/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer *grpc.Server
	pvzService *grpcservice.PVZService
	config     *config.Config
}

func NewServer(conf *config.Config, pvzService *grpcservice.PVZService) *Server {
	grpcServer := grpc.NewServer()

	pvz_v1.RegisterPVZServiceServer(grpcServer, pvzService)

	reflection.Register(grpcServer)

	return &Server{
		grpcServer: grpcServer,
		pvzService: pvzService,
		config:     conf,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.config.GRPC.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	log.Printf("Starting gRPC server on %s", addr)

	if err := s.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
