package server

import (
	"STTMain/internal/service"
	"context"
	"net"

	sttv1 "github.com/skinkvi/protosSTT/gen/go/stt"
	"google.golang.org/grpc"
)

type Server struct {
	sttService service.SpeedTypingTestService
	grpcServer *grpc.Server
}

func NewServer(sttService service.SpeedTypingTestService) *Server {
	grpcServer := grpc.NewServer()
	sttv1.RegisterSpeedTypingTestServer(grpcServer, &sttService)

	return &Server{
		sttService: sttService,
		grpcServer: grpcServer,
	}
}

func (s *Server) Start(listener net.Listener) error {
	return s.grpcServer.Serve(listener)
}

func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.GracefulStop()
	return nil
}
