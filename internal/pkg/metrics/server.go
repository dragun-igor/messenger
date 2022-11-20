package metrics

import (
	"net/http"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

type MetricsServerService struct {
	httpServer        *http.Server
	grpcServerMetrics *grpc_prometheus.ServerMetrics
}

func NewMetricsServerService(addr string) *MetricsServerService {
	reg := prometheus.NewRegistry()
	grpcServerMetrics := grpc_prometheus.NewServerMetrics()
	reg.MustRegister(grpcServerMetrics)
	return &MetricsServerService{
		httpServer:        &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: addr},
		grpcServerMetrics: grpcServerMetrics,
	}
}

func (s *MetricsServerService) Initialize(server *grpc.Server) {
	s.grpcServerMetrics.InitializeMetrics(server)
}

func (s *MetricsServerService) GRPCServerUnaryMetricsInterceptor() grpc.UnaryServerInterceptor {
	return s.grpcServerMetrics.UnaryServerInterceptor()
}

func (s *MetricsServerService) GRPCServerStreamMetricsInterceptor() grpc.StreamServerInterceptor {
	return s.grpcServerMetrics.StreamServerInterceptor()
}

func (s *MetricsServerService) Listen() error {
	return s.httpServer.ListenAndServe()
}

func (s *MetricsServerService) Close() error {
	return s.httpServer.Close()
}
