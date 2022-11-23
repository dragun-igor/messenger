package client

import (
	"context"
	"log"

	"github.com/dragun-igor/messenger/internal/client/service"
	"github.com/dragun-igor/messenger/internal/pkg/metrics"
	"github.com/dragun-igor/messenger/proto/messenger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	service *service.Service
	metrics *metrics.MetricsClientService
	conn    *grpc.ClientConn
}

func New(grpcAddr, promAddr string) (*Client, error) {
	client := &Client{}
	metric := metrics.NewMetricsClientService(promAddr)
	conn, err := grpc.Dial(grpcAddr,
		grpc.WithChainUnaryInterceptor(
			metric.GRPCClientUnaryMetricsInterceptor(),
			metric.AppMetricsInterceptor(),
		),
		grpc.WithStreamInterceptor(metric.GRPCClientStreamMetricsInterceptor()),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client.metrics = metric
	client.conn = conn
	return client, nil
}

func (c *Client) Serve(ctx context.Context) error {
	defer c.Stop()
	c.service = service.New(messenger.NewMessengerServiceClient(c.conn))
	go func() {
		if err := c.metrics.Listen(); err != nil {
			log.Println(err)
		}
	}()
	return c.service.Serve(ctx)
}

func (c *Client) Stop() {
	log.Println("closing connection")
	c.conn.Close()

	log.Println("stop metrics server")
	c.metrics.Close()
}
