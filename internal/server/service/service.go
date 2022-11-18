package service

import (
	"log"

	"github.com/dragun-igor/messenger/config"
	"github.com/dragun-igor/messenger/internal/server/model"
	"github.com/dragun-igor/messenger/internal/server/resources"
	"github.com/dragun-igor/messenger/proto/messenger"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MessengerServiceServer struct {
	messenger.UnimplementedMessengerServiceServer
	clients   map[string]chan *messenger.Message
	resources *resources.Resources
}

func NewMessengerServiceServer(ctx context.Context, config *config.Config) (*MessengerServiceServer, error) {
	resources, err := resources.NewResources(ctx, "migrations/createtables.up.sql", config)
	if err != nil {
		return nil, err
	}
	return &MessengerServiceServer{
		resources: resources,
		clients:   make(map[string]chan *messenger.Message),
	}, nil
}

func (s *MessengerServiceServer) SendMessage(ctx context.Context, message *messenger.Message) (*messenger.MessageResponse, error) {
	if _, ok := s.clients[message.Receiver]; !ok {
		return &messenger.MessageResponse{Sent: false}, nil
	}
	msg := model.Message{
		Sender:   message.Sender,
		Receiver: message.Receiver,
		Message:  message.Message,
	}
	s.resources.InsertMessage(ctx, msg)
	s.clients[message.Receiver] <- message
	log.Printf("%s -> %s: %s", message.Sender, message.Receiver, message.Message)
	return &messenger.MessageResponse{Sent: true}, nil
}

func (s *MessengerServiceServer) ReceiveMessage(user *messenger.User, stream messenger.MessengerService_ReceiveMessageServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case msg := <-s.clients[user.Name]:
			stream.Send(msg)
		}
	}
}

func (s *MessengerServiceServer) Ping(stream messenger.MessengerService_PingServer) error {
	user, err := stream.Recv()
	if err != nil {
		return err
	}
	log.Printf("user %s is online", user.Name)
	s.clients[user.Name] = make(chan *messenger.Message)
	if err := stream.Send(&emptypb.Empty{}); err != nil {
		return err
	}
	<-stream.Context().Done()
	delete(s.clients, user.Name)
	log.Printf("user %s is offline", user.Name)
	return nil
}
