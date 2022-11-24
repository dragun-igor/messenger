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
	"github.com/dragun-igor/messenger/internal/server/resources"
	"github.com/dragun-igor/messenger/internal/server/service"
	"github.com/dragun-igor/messenger/proto/messenger"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

const gracefulTimeout = 2 * time.Second

type Server struct {
	grpc    *grpc.Server
	db      resources.PostgresDB
	config  *config.Config
	metrics *metrics.ServerMetrics
	closeCh chan struct{}
}

func New(ctx context.Context, config *config.Config) (*Server, error) {
	server := &Server{}
	db, err := resources.InitPostgresDB(ctx, config)
	if err != nil {
		return nil, err
	}
	server.closeCh = make(chan struct{})
	server.config = config
	server.db = db
	server.metrics = metrics.NewServerMetrics(config.PrometheusHost + ":" + config.PrometheusPort)
	server.grpc = grpc.NewServer(
		grpc.StreamInterceptor(server.metrics.GRPCServerStreamMetricsInterceptor()),
		grpc.ChainUnaryInterceptor(
			server.metrics.GRPCServerUnaryMetricsInterceptor(),
			server.metrics.AppMetricsInterceptor(),
		),
	)
	messenger.RegisterMessengerServiceServer(server.grpc, service.NewServiceServer(server.db, server.closeCh))
	err = server.metrics.Initialize(server.grpc)
	if err != nil {
		return nil, err
	}
	grpc_prometheus.Register(server.grpc)
	return server, nil
}

func (s *Server) Serve(ctx context.Context) error {
	defer s.Stop(ctx)
	lis, err := net.Listen("tcp", s.config.GRPCHost+":"+s.config.GRPCPort)
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
