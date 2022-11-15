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

func (s *messengerServiceServer) SignUp(ctx context.Context, signUpData *messengerpb.SignUpData) (*messengerpb.UserData, error) {
	userData, err := s.resources.SignUp(ctx, signUpData)
	if err != nil {
		return nil, err
	}
	log.Printf("User %s ID %d logged in\n", userData.Name, userData.Id)
	return userData, nil
}

func (s *messengerServiceServer) SignIn(ctx context.Context, signInData *messengerpb.SignInData) (*messengerpb.UserData, error) {
	userData, err := s.resources.SignIn(ctx, signInData)
	if err != nil {
		return nil, err
	}
	log.Printf("User %s ID %d logged in\n", userData.Name, userData.Id)
	return userData, nil
}

func (s *messengerServiceServer) CheckName(ctx context.Context, checkNameMessage *messengerpb.CheckNameMessage) (*messengerpb.CheckNameAck, error) {
	ack, err := s.resources.CheckName(ctx, checkNameMessage)
	if err != nil {
		fmt.Printf("err: %v \n", err)
	}
	fmt.Println(ack.Busy)
	fmt.Println(err)
	return ack, err
}

func (s *messengerServiceServer) CheckLogin(ctx context.Context, checkLoginMessage *messengerpb.CheckLoginMessage) (*messengerpb.CheckLoginAck, error) {
	ack, err := s.resources.CheckLogin(ctx, checkLoginMessage)
	if err != nil {
		fmt.Printf("err: %v \n", err)
	}
	return ack, err
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
		for !s.resources.SendMessage(context.Background(), msg) {
		}
		s.clients[msg.Receiver.Id] <- msg
	}()
	return nil
}

func (s *messengerServiceServer) ReceiveMessage(userData *messengerpb.UserData, msgStream messengerpb.MessengerService_ReceiveMessageServer) error {
	msgCh := make(chan *messengerpb.Message)
	s.clients[userData.Id] = msgCh
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
	// ctx, cancel := context.WithCancel(context.Background())
	// go func() {
	// 	exit := make(chan os.Signal, 1)
	// 	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	// 	cancel()
	// }()
	return &messengerServiceServer{
		clients:   make(map[int64]chan *messengerpb.Message),
		resources: resources.GetResources(context.Background(), config.New()),
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
