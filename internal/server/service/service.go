package service

import (
	"log"

	"github.com/dragun-igor/messenger/internal/pkg/model"
	"github.com/dragun-igor/messenger/pkg/errors"
	"github.com/dragun-igor/messenger/proto/messenger"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	messenger.UnimplementedMessengerServiceServer
	clients map[string]chan *messenger.Message
	db      Repository
	closeCh <-chan struct{}
}

func New(db Repository, closeCh <-chan struct{}) *Service {
	return &Service{
		db:      db,
		clients: make(map[string]chan *messenger.Message),
		closeCh: closeCh,
	}
}

func (s *Service) SendMessage(ctx context.Context, message *messenger.Message) (*messenger.MessageResponse, error) {
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

func (s *Service) Ping(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *Service) ReceiveMessage(stream messenger.MessengerService_ReceiveMessageServer) error {
	user, err := stream.Recv()
	if err != nil {
		return convert(err)
	}
	s.clients[user.Name] = make(chan *messenger.Message)
	log.Printf("user %s is online", user.Name)
	for {
		select {
		case <-s.closeCh:
			delete(s.clients, user.Name)
			return nil
		case <-stream.Context().Done():
			log.Printf("user %s is offline", user.Name)
			delete(s.clients, user.Name)
			return nil
		case msg := <-s.clients[user.Name]:
			err := stream.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (s *Service) SignUp(ctx context.Context, signUpRequest *messenger.SignUpRequest) (*emptypb.Empty, error) {
	user := model.AuthData{
		Login: signUpRequest.Login,
		Name:  signUpRequest.Name,
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
	if err := user.SetHashByPassword(signUpRequest.Password); err != nil {
		return nil, convert(err)
	}
	err = s.db.CreateUser(ctx, user)
	if err != nil {
		return nil, convert(err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) LogIn(ctx context.Context, logInRequest *messenger.LogInRequest) (*messenger.User, error) {
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
