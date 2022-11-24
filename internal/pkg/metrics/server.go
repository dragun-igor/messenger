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

type ServerMetrics struct { //nolint:revive
	httpServer        *http.Server
	grpcServerMetrics *grpc_prometheus.ServerMetrics
}

func NewServerMetrics(addr string) *ServerMetrics {
	reg := prometheus.NewRegistry()
	grpcServerMetrics := grpc_prometheus.NewServerMetrics()
	reg.MustRegister(grpcServerMetrics)
	reg.MustRegister(requestTimeHist)
	reg.MustRegister(requestErrorsCounter)
	return &ServerMetrics{
		httpServer: &http.Server{
			Handler:           promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
			Addr:              addr,
			ReadHeaderTimeout: 3 * time.Second,
		},
		grpcServerMetrics: grpcServerMetrics,
	}
}

func (s *ServerMetrics) Initialize(server *grpc.Server) error {
	s.grpcServerMetrics.InitializeMetrics(server)
	_, err := requestTimeHist.GetMetricWithLabelValues(fieldMethodName)
	if err != nil {
		return err
	}
	_, err = requestErrorsCounter.GetMetricWithLabelValues(fieldMethodName)
	return err
}

func (s *ServerMetrics) GRPCServerUnaryMetricsInterceptor() grpc.UnaryServerInterceptor {
	return s.grpcServerMetrics.UnaryServerInterceptor()
}

func (s *ServerMetrics) GRPCServerStreamMetricsInterceptor() grpc.StreamServerInterceptor {
	return s.grpcServerMetrics.StreamServerInterceptor()
}

func (s *ServerMetrics) AppMetricsInterceptor() grpc.UnaryServerInterceptor {
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

func (s *ServerMetrics) Listen() error {
	return s.httpServer.ListenAndServe()
}

func (s *ServerMetrics) Close() error {
	return s.httpServer.Close()
}
