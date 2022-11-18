package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/dragun-igor/messenger/proto/messenger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client messenger.MessengerServiceClient
	server string
	name   string
}

func NewClient(server, name string) *Client {
	return &Client{
		server: server,
		name:   name,
	}
}

func (c *Client) Serve() error {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(c.server, opts...)
	if err != nil {
		return err
	}
	defer conn.Close()
	ctx := context.Background()
	c.client = messenger.NewMessengerServiceClient(conn)
	if err := c.ping(ctx); err != nil {
		return err
	}
	if err := c.listenScanner(ctx); err != nil {
		return err
	}
	return nil
}

func (c *Client) listenScanner(ctx context.Context) error {
	fmt.Printf("[SERVICE] Hello, %s!\n", c.name)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		switch message {
		case "":
			continue
		default:
			if err := c.sendMessage(ctx, message); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Client) receiveMessage(ctx context.Context) {
	stream, err := c.client.ReceiveMessage(ctx, &messenger.User{Name: c.name})
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		<-stream.Context().Done()
		fmt.Println("[SERVICE] Connection to server has lost")
	}()
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			return
		}
		fmt.Printf("%v: %v\n", msg.Sender, msg.Message)
	}
}

func (c *Client) sendMessage(ctx context.Context, message string) error {
	messageSplit := strings.SplitN(message, " ", 2)
	response, err := c.client.SendMessage(ctx, &messenger.Message{
		Sender:   c.name,
		Receiver: messageSplit[0],
		Message:  messageSplit[1],
	})
	if err != nil {
		return err
	}
	if !response.Sent {
		fmt.Printf("[SERVICE] User %s is offline!\n", messageSplit[0])
	}
	return nil
}

func (c *Client) ping(ctx context.Context) error {
	stream, err := c.client.Ping(ctx)
	if err != nil {
		return err
	}
	err = stream.Send(&messenger.User{Name: c.name})
	if err != nil {
		return err
	}
	_, err = stream.Recv()
	if err != nil {
		return err
	}
	go c.receiveMessage(ctx)
	return nil
}
