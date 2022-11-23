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
	service *service.Service
	db      resources.PostgresDB
	config  *config.Config
	metrics *metrics.MetricsServerService
}

func New(ctx context.Context, config *config.Config) (*Server, error) {
	server := &Server{}
	db, err := resources.InitPostgresDB(ctx, config)
	if err != nil {
		return nil, err
	}
	s := service.New(ctx, server.db)
	server.config = config
	server.db = db
	server.service = s
	server.metrics = metrics.NewMetricsServerService(config.PrometheusHost + ":" + config.PrometheusPort)
	server.grpc = grpc.NewServer(
		grpc.StreamInterceptor(server.metrics.GRPCServerStreamMetricsInterceptor()),
		grpc.ChainUnaryInterceptor(
			server.metrics.GRPCServerUnaryMetricsInterceptor(),
			server.metrics.AppMetricsInterceptor(),
		),
	)
	messenger.RegisterMessengerServiceServer(server.grpc, s)
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

	log.Println("stopping grpc server")
	s.grpc.Stop()

	log.Println("disconnecting db")
	s.db.Close(context.Background())

	log.Println("stopping metrics server")
	s.metrics.Close()
}
