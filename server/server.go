package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/dragun-igor/messenger/config"
	"github.com/dragun-igor/messenger/messengerpb"
	"github.com/dragun-igor/messenger/resources"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type messengerServiceServer struct {
	messengerpb.UnimplementedMessengerServiceServer
	clients   map[int64]chan *messengerpb.Message
	resources *resources.Resources
}

func (s *messengerServiceServer) SignIn(context context.Context, signInData *messengerpb.SignInData) (*messengerpb.User, error) {
	id, name, err := s.resources.SignIn(signInData)
	if err != nil {
		return nil, err
	}
	log.Printf("User %s ID %d logged in\n", name, id)
	return &messengerpb.User{Id: id, FirstName: name}, nil
}

func (s *messengerServiceServer) SendMessage(msgStream messengerpb.MessengerService_SendMessageServer) error {
	msg, err := msgStream.Recv()
	msg.Time = timestamppb.Now()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	ack := &messengerpb.MessageAck{Status: "received"}
	msgStream.SendAndClose(ack)
	go func() {
		for !s.resources.SendMessage(msg) {
		}
		s.clients[msg.Receiver.Id] <- msg
	}()
	return nil
}

func (s *messengerServiceServer) ReceiveMessage(userID *messengerpb.User, msgStream messengerpb.MessengerService_ReceiveMessageServer) error {
	msgCh := make(chan *messengerpb.Message)
	s.clients[userID.Id] = msgCh
	for {
		select {
		case <-msgStream.Context().Done():
			return nil
		case msg := <-msgCh:
			fmt.Printf("%v -> %v: %v \n", msg.Sender.Id, msg.Receiver.Id, msg.Message)
			msgStream.Send(msg)
		}
	}
}

func newServer() *messengerServiceServer {
	context := context.Background()
	return &messengerServiceServer{
		clients:   make(map[int64]chan *messengerpb.Message),
		resources: resources.GetResources(context, config.New()),
	}
}

func main() {
	fmt.Println("--- SERVER APP ---")
	lis, err := net.Listen("tcp", ":5400")
	if err != nil {
		log.Fatalf("Failed to listen: %v \n", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	messengerpb.RegisterMessengerServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
