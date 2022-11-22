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
	metrics *metrics.MetricsServerService
}

func NewServer(ctx context.Context, config *config.Config) (*Server, error) {
	server := &Server{}
	db, err := resources.InitPostgresDB(ctx, config)
	if err != nil {
		return nil, err
	}
	server.config = config
	server.db = db
	server.metrics = metrics.NewMetricsServerService(config.PrometheusHost + ":" + config.PrometheusPort)
	server.grpc = grpc.NewServer(
		grpc.StreamInterceptor(server.metrics.GRPCServerStreamMetricsInterceptor()),
		grpc.ChainUnaryInterceptor(
			server.metrics.GRPCServerUnaryMetricsInterceptor(),
			server.metrics.AppMetricsInterceptor(),
		),
	)
	messenger.RegisterMessengerServiceServer(server.grpc, service.NewMessengerServiceServer(ctx, server.db))
	server.metrics.Initialize(server.grpc)
	grpc_prometheus.Register(server.grpc)
	return server, nil
}

func (s *Server) Serve() error {
	defer s.Stop()
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

func (s *Server) Stop() {
	time.Sleep(gracefulTimeout)

	log.Println("disconnecting db")
	s.db.Close(context.Background())
	log.Println("connection to db has closed")

	log.Println("stopping metrics server")
	s.metrics.Close()
	log.Println("metrics server is stopped")
}
