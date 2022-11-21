package client

import (
	"context"
	"log"
	"time"

	"github.com/dragun-igor/messenger/internal/client/service"
	"github.com/dragun-igor/messenger/internal/pkg/metrics"
	"github.com/dragun-igor/messenger/proto/messenger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const gracefulTimeout = 2 * time.Second

type Client struct {
	service *service.MessengerServiceClient
	metrics *metrics.MetricsClientService
}

func NewClient(phost, pport string) *Client {
	return &Client{
		metrics: metrics.NewMetricsClientService(phost + ":" + pport),
	}
}

func (c *Client) Serve(ghost, gport string) error {
	defer c.Stop()
	conn, err := grpc.Dial(ghost+":"+gport,
		grpc.WithChainUnaryInterceptor(
			c.metrics.GRPCClientUnaryMetricsInterceptor(),
			c.metrics.AppMetricsInterceptor(),
		),
		grpc.WithStreamInterceptor(c.metrics.GRPCClientStreamMetricsInterceptor()),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx := context.Background()
	c.service = service.NewClientService(messenger.NewMessengerServiceClient(conn))

	go func() {
		if err := c.metrics.Listen(); err != nil {
			log.Println(err)
		}
	}()

	return c.service.Serve(ctx)
}

func (c *Client) Stop() {
	time.Sleep(gracefulTimeout)

	log.Println("stop metrics server")
	c.metrics.Close()
}
