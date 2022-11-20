package metrics

import (
	"net/http"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

type MetricsClientService struct {
	httpServer        *http.Server
	grpcClientMetrics *grpc_prometheus.ClientMetrics
}

func NewMetricsClientService(addr string) *MetricsClientService {
	reg := prometheus.NewRegistry()
	grpcClientMetrics := grpc_prometheus.NewClientMetrics()
	reg.MustRegister(grpcClientMetrics)
	return &MetricsClientService{
		httpServer:        &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: addr},
		grpcClientMetrics: grpcClientMetrics,
	}
}

func (c *MetricsClientService) GRPCClientUnaryMetricsInterceptor() grpc.UnaryClientInterceptor {
	return c.grpcClientMetrics.UnaryClientInterceptor()
}

func (c *MetricsClientService) GRPCClientStreamMetricsInterceptor() grpc.StreamClientInterceptor {
	return c.grpcClientMetrics.StreamClientInterceptor()
}

func (c *MetricsClientService) Listen() error {
	return c.httpServer.ListenAndServe()
}

func (c *MetricsClientService) Close() error {
	return c.httpServer.Close()
}
