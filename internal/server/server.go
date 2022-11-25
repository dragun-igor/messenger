package server

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dragun-igor/messenger/config"
	"github.com/dragun-igor/messenger/internal/pkg/metrics"
	"github.com/dragun-igor/messenger/internal/pkg/repository"
	"github.com/dragun-igor/messenger/internal/server/resources"
	"github.com/dragun-igor/messenger/internal/server/service"
	"github.com/dragun-igor/messenger/proto/messenger"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

const gracefulTimeout = 2 * time.Second

type Server struct {
	grpc    *grpc.Server
	db      resources.Connection
	config  *config.Config
	metrics *metrics.ServerMetrics
	closeCh chan struct{}
}

func New(config *config.Config) (*Server, error) {
	ctx := context.TODO()
	server := &Server{}
	db, err := resources.NewConnection(ctx, config)
	if err != nil {
		return nil, err
	}
	server.closeCh = make(chan struct{})
	server.config = config
	server.db = db
	server.metrics = metrics.NewServerMetrics(config.PrometheusAddr)
	server.grpc = grpc.NewServer(
		grpc.StreamInterceptor(server.metrics.GRPCServerStreamMetricsInterceptor()),
		grpc.ChainUnaryInterceptor(
			server.metrics.GRPCServerUnaryMetricsInterceptor(),
			server.metrics.AppMetricsInterceptor(),
		),
	)
	messenger.RegisterMessengerServer(server.grpc, service.NewServiceServer(repository.New(db), server.closeCh))
	err = server.metrics.Initialize(server.grpc)
	if err != nil {
		return nil, err
	}
	grpc_prometheus.Register(server.grpc)
	return server, nil
}

func (s *Server) Serve(ctx context.Context) error {
	defer s.Stop(ctx)
	lis, err := net.Listen("tcp", s.config.GRPCAddr)
	if err != nil {
		return err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("termination signal received")
		log.Println("stopping grpc server")
		close(s.closeCh)
		s.grpc.GracefulStop()
		log.Println("grpc server is stopped")
	}()

	go func() {
		if err := s.metrics.Listen(); err != nil {
			log.Println(err)
		}
	}()
	return s.grpc.Serve(lis)
}

func (s *Server) Stop(ctx context.Context) {
	time.Sleep(gracefulTimeout)

	log.Println("disconnecting db")
	s.db.Close(ctx)

	log.Println("stopping metrics server")
	s.metrics.Close()
}
