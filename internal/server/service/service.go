package service

import (
	"log"

	"github.com/dragun-igor/messenger/internal/server/model"
	"github.com/dragun-igor/messenger/pkg/errors"
	"github.com/dragun-igor/messenger/proto/messenger"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MessengerServiceServer struct {
	messenger.UnimplementedMessengerServiceServer
	clients map[string]chan *messenger.Message
	db      Repository
}

func NewMessengerServiceServer(ctx context.Context, db Repository) *MessengerServiceServer {
	return &MessengerServiceServer{
		db:      db,
		clients: make(map[string]chan *messenger.Message),
	}
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
	if err := s.db.InsertMessage(ctx, msg); err != nil {
		return nil, convert(err)
	}
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
		return convert(err)
	}
	log.Printf("user %s is online", user.Name)
	s.clients[user.Name] = make(chan *messenger.Message)
	if err := stream.Send(&emptypb.Empty{}); err != nil {
		return convert(err)
	}
	<-stream.Context().Done()
	delete(s.clients, user.Name)
	log.Printf("user %s is offline", user.Name)
	return nil
}

func (s *MessengerServiceServer) SignUp(ctx context.Context, signUpRequest *messenger.SignUpRequest) (*emptypb.Empty, error) {
	user := model.AuthData{
		Login: signUpRequest.Login,
		Name:  signUpRequest.Name,
	}
	if err := user.SetHashByPassword(signUpRequest.Password); err != nil {
		return nil, convert(err)
	}
	ok, err := s.db.CheckLoginExists(ctx, user)
	if err != nil {
		return nil, convert(err)
	}
	if !ok {
		return nil, convert(errors.ErrLoginNameIsBusy)
	}
	ok, err = s.db.CheckNameExists(ctx, user)
	if err != nil {
		return nil, convert(err)
	}
	if !ok {
		return nil, convert(errors.ErrUserNameIsBusy)
	}
	err = s.db.CreateUser(ctx, user)
	if err != nil {
		return nil, convert(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *MessengerServiceServer) LogIn(ctx context.Context, logInRequest *messenger.LogInRequest) (*messenger.User, error) {
	user := model.AuthData{
		Login: logInRequest.Login,
	}
	dbUser, err := s.db.GetUser(ctx, user)
	if err != nil {
		return nil, convert(err)
	}
	if !dbUser.IsPasswordCorrect(logInRequest.Password) {
		return nil, convert(errors.ErrIncorrectPassword)
	}
	return &messenger.User{Name: dbUser.Name}, nil
}
