package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Server struct {
	port string
	srv  *grpc.Server
}

func New(port string) *Server {
	return &Server{
		port: port,
		srv:  grpc.NewServer(),
	}
}

func (s *Server) Server() *grpc.Server {
	return s.srv
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	log.Info().Msgf("gRPC server listening on :%s", s.port)
	return s.srv.Serve(listener)
}

func (s *Server) Stop(_ context.Context) error {
	log.Info().Msg("gRPC server stopping...")
	s.srv.GracefulStop()
	return nil
}
