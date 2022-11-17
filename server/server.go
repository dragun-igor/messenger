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

type Client struct {
	Message          chan *messengerpb.Message
	RequestToFriends chan *messengerpb.UsersPair
	User             chan *messengerpb.User
}

type messengerServiceServer struct {
	messengerpb.UnimplementedMessengerServiceServer
	clients   map[string]*Client
	resources *resources.Resources
}

func (s *messengerServiceServer) SignUp(ctx context.Context, signUpData *messengerpb.SignUpData) (*messengerpb.User, error) {
	user, err := s.resources.SignUp(ctx, signUpData)
	if err != nil {
		return nil, err
	}
	s.clients[user.Name] = &Client{
		Message:          make(chan *messengerpb.Message),
		RequestToFriends: make(chan *messengerpb.UsersPair),
	}
	log.Printf("user %s signed up and logged in\n", user.Name)
	return user, nil
}

func (s *messengerServiceServer) LogIn(ctx context.Context, logInData *messengerpb.LogInData) (*messengerpb.User, error) {
	user, err := s.resources.SignIn(ctx, logInData)
	if err != nil {
		log.Printf("err: %v \n", err)
		return nil, err
	}
	log.Printf("user %s logged in\n", user.Name)
	return user, nil
}

func (s *messengerServiceServer) CheckName(ctx context.Context, checkNameMessage *messengerpb.CheckNameMessage) (*messengerpb.MessageAck, error) {
	ack, err := s.resources.CheckName(ctx, checkNameMessage)
	if err != nil {
		log.Printf("err: %v \n", err)
	}
	return ack, err
}

func (s *messengerServiceServer) CheckLogin(ctx context.Context, checkLoginMessage *messengerpb.CheckLoginMessage) (*messengerpb.MessageAck, error) {
	ack, err := s.resources.CheckLogin(ctx, checkLoginMessage)
	if err != nil {
		log.Printf("err: %v \n", err)
	}
	return ack, err
}

func (s *messengerServiceServer) RequestAddToFriendsList(ctx context.Context, usersPair *messengerpb.UsersPair) (*messengerpb.MessageAck, error) {
	ack, err := s.resources.RequestAddToFriendsList(ctx, usersPair)
	if err != nil {
		log.Printf("err: %v \n", err)
		return nil, err
	}
	go func() {
		s.clients[usersPair.Receiver].RequestToFriends <- usersPair
	}()
	return ack, nil
}

func (s *messengerServiceServer) ListenAddToFriendsList(user *messengerpb.User, userStream messengerpb.MessengerService_ListenAddToFriendsListServer) error {
	for {
		select {
		case <-userStream.Context().Done():
			return nil
		case request := <-s.clients[user.Name].RequestToFriends:
			log.Printf("user %v requested adding to friends list user %v \n", request.Requester, request.Receiver)
			userStream.Send(&messengerpb.User{Name: request.Requester})
		}
	}
}

func (s *messengerServiceServer) ListenAppendNewFriend(user *messengerpb.User, userStream messengerpb.MessengerService_ListenAppendNewFriendServer) error {
	for {
		select {
		case <-userStream.Context().Done():
			return nil
		case user := <-s.clients[user.Name].User:
			userStream.Send(user)
		}
	}
}

func (s *messengerServiceServer) AddToFriendsList(ctx context.Context, usersPair *messengerpb.UsersPair) (*messengerpb.MessageAck, error) {
	ack, err := s.resources.AddToFriendsList(ctx, usersPair)
	go func() {
		s.clients[usersPair.Requester].User <- &messengerpb.User{Name: usersPair.Receiver}
	}()
	return ack, err
}

func (s *messengerServiceServer) GetMessages(ctx context.Context, usersPair *messengerpb.UsersPair) (*messengerpb.MessageArchive, error) {
	archive, err := s.resources.GetMessages(ctx, usersPair)
	if err != nil {
		log.Printf("err: %v \n", err)
	}
	return archive, err
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
	ack := &messengerpb.MessageAck{Status: "Message sent"}
	msgStream.SendAndClose(ack)
	err = s.resources.SendMessage(context.Background(), msg)
	if err != nil {
		return err
	}
	log.Printf("%v -> %v: %v \n", msg.Sender, msg.Receiver, msg.Message)
	s.clients[msg.Receiver].Message <- msg
	return nil
}

func (s *messengerServiceServer) ReceiveMessage(user *messengerpb.User, msgStream messengerpb.MessengerService_ReceiveMessageServer) error {
	for {
		select {
		case <-msgStream.Context().Done():
			return nil
		case msg := <-s.clients[user.Name].Message:
			msgStream.Send(msg)
		}
	}
}

func (s *messengerServiceServer) GetFriendsList(ctx context.Context, user *messengerpb.User) (*messengerpb.FriendsList, error) {
	friends, err := s.resources.GetFriendsList(ctx, user)
	if err != nil {
		return nil, err
	}
	return friends, nil
}

func (s *messengerServiceServer) GetRequestsToFriendsList(ctx context.Context, user *messengerpb.User) (*messengerpb.Requests, error) {
	requests, err := s.resources.GetRequestsToFriendsList(ctx, user)
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func newServer() *messengerServiceServer {
	log.Println("connecting to database")
	resources := resources.GetResources(context.Background(), config.New())
	log.Println("getting users list")
	users, err := resources.GetAllUsers(context.Background())
	if err != nil {
		log.Fatalf("err: %v \n", err)
	}
	log.Println("init users channels")
	clients := make(map[string]*Client, len(users))
	for _, name := range users {
		clients[name] = &Client{
			Message:          make(chan *messengerpb.Message),
			RequestToFriends: make(chan *messengerpb.UsersPair),
			User:             make(chan *messengerpb.User),
		}
	}
	log.Println("succesfully created server")
	return &messengerServiceServer{
		clients:   clients,
		resources: resources,
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
