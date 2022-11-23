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

type MetricsClientService struct { //nolint:revive
	httpServer        *http.Server
	grpcClientMetrics *grpc_prometheus.ClientMetrics
}

func NewMetricsClientService(addr string) *MetricsClientService {
	reg := prometheus.NewRegistry()
	grpcClientMetrics := grpc_prometheus.NewClientMetrics()
	reg.MustRegister(grpcClientMetrics)
	reg.MustRegister(requestTimeHist)
	reg.MustRegister(requestErrorsCounter)
	return &MetricsClientService{
		httpServer: &http.Server{
			Handler:           promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
			Addr:              addr,
			ReadHeaderTimeout: 3 * time.Second,
		},
		grpcClientMetrics: grpcClientMetrics,
	}
}

func (c *MetricsClientService) GRPCClientUnaryMetricsInterceptor() grpc.UnaryClientInterceptor {
	return c.grpcClientMetrics.UnaryClientInterceptor()
}

func (c *MetricsClientService) GRPCClientStreamMetricsInterceptor() grpc.StreamClientInterceptor {
	return c.grpcClientMetrics.StreamClientInterceptor()
}

func (c *MetricsClientService) AppMetricsInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		requestTimeHist.WithLabelValues(method).Observe(time.Since(start).Seconds())
		if err != nil {
			requestErrorsCounter.WithLabelValues(method).Inc()
		}
		return err
	}
}

func (c *MetricsClientService) Listen() error {
	return c.httpServer.ListenAndServe()
}

func (c *MetricsClientService) Close() error {
	return c.httpServer.Close()
}
