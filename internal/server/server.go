package server

import (
	"context"
	"fmt"
	"net"

	"github.com/dragun-igor/messenger/config"
	"github.com/dragun-igor/messenger/internal/server/service"
	"github.com/dragun-igor/messenger/proto/messenger"
	"google.golang.org/grpc"
)

type Server struct {
	config *config.Config
}

func NewServer(ctx context.Context, config *config.Config) (*Server, error) {
	return &Server{
		config: config,
	}, nil
}

func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.config.GRPCHost, s.config.GRPCPort))
	if err != nil {
		return err
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	serv, err := service.NewMessengerServiceServer(context.Background(), s.config)
	if err != nil {
		return err
	}
	messenger.RegisterMessengerServiceServer(grpcServer, serv)
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
