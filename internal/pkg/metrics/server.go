package metrics

import (
	"context"
	"net/http"
	"time"

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
	reg.MustRegister(requestTimeHist)
	reg.MustRegister(requestErrorsCounter)
	return &MetricsServerService{
		httpServer:        &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: addr},
		grpcServerMetrics: grpcServerMetrics,
	}
}

func (s *MetricsServerService) Initialize(server *grpc.Server) error {
	s.grpcServerMetrics.InitializeMetrics(server)
	_, err := requestTimeHist.GetMetricWithLabelValues(fieldMethodName)
	if err != nil {
		return err
	}
	_, err = requestErrorsCounter.GetMetricWithLabelValues(fieldMethodName)
	return err
}

func (s *MetricsServerService) GRPCServerUnaryMetricsInterceptor() grpc.UnaryServerInterceptor {
	return s.grpcServerMetrics.UnaryServerInterceptor()
}

func (s *MetricsServerService) GRPCServerStreamMetricsInterceptor() grpc.StreamServerInterceptor {
	return s.grpcServerMetrics.StreamServerInterceptor()
}

func (s *MetricsServerService) AppMetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		requestTimeHist.WithLabelValues(info.FullMethod).Observe(time.Since(start).Seconds())
		if err != nil {
			requestErrorsCounter.WithLabelValues(info.FullMethod).Inc()
		}
		return resp, err
	}
}

func (s *MetricsServerService) Listen() error {
	return s.httpServer.ListenAndServe()
}

func (s *MetricsServerService) Close() error {
	return s.httpServer.Close()
}
