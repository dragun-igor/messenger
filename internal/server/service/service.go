package service

import (
	"log"

	"github.com/dragun-igor/messenger/config"
	"github.com/dragun-igor/messenger/internal/server/model"
	"github.com/dragun-igor/messenger/internal/server/resources"
	"github.com/dragun-igor/messenger/pkg/errors"
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

func (s *MessengerServiceServer) SignUp(ctx context.Context, signUpRequest *messenger.SignUpRequest) (*emptypb.Empty, error) {
	user := model.User{
		Login: signUpRequest.Login,
		Name:  signUpRequest.Name,
	}
	user.SetHashByPassword(signUpRequest.Password)
	ok, err := s.resources.CheckLoginExists(ctx, user)
	if err != nil {
		return nil, convert(err)
	}
	if !ok {
		return nil, convert(errors.ErrLoginNameIsBusy)
	}
	ok, err = s.resources.CheckNameExists(ctx, user)
	if err != nil {
		return nil, convert(err)
	}
	if !ok {
		return nil, convert(errors.ErrUserNameIsBusy)
	}
	err = s.resources.InsertUser(ctx, user)
	if err != nil {
		return nil, convert(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *MessengerServiceServer) LogIn(ctx context.Context, logInRequest *messenger.LogInRequest) (*messenger.User, error) {
	user := model.User{
		Login:    logInRequest.Login,
		Password: logInRequest.Password,
	}
	name, err := s.resources.LogIn(ctx, user)
	if err != nil {
		return nil, convert(err)
	}
	return &messenger.User{Name: name}, nil
}
