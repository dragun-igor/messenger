package service

import (
	"log"
	"sync"

	"github.com/dragun-igor/messenger/internal/pkg/model"
	"github.com/dragun-igor/messenger/pkg/errors"
	"github.com/dragun-igor/messenger/proto/messenger"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ServiceServer struct {
	messenger.MessengerServer
	mu      sync.Mutex
	clients map[string]chan *messenger.Message
	db      Repository
	closeCh <-chan struct{}
}

func NewServiceServer(db Repository, closeCh <-chan struct{}) *ServiceServer {
	return &ServiceServer{
		mu:      sync.Mutex{},
		db:      db,
		clients: make(map[string]chan *messenger.Message),
		closeCh: closeCh,
	}
}

func (s *ServiceServer) SendMessage(ctx context.Context, message *messenger.Message) (*messenger.MessageResponse, error) {
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

func (s *ServiceServer) Ping(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *ServiceServer) ReceiveMessage(stream messenger.Messenger_ReceiveMessageServer) error {
	user, err := stream.Recv()
	if err != nil {
		return convert(err)
	}
	s.mu.Lock()
	s.clients[user.Name] = make(chan *messenger.Message)
	s.mu.Unlock()
	log.Printf("user %s is online", user.Name)
	for {
		select {
		case <-s.closeCh:
			return nil
		case <-stream.Context().Done():
			log.Printf("user %s is offline", user.Name)
			s.mu.Lock()
			delete(s.clients, user.Name)
			s.mu.Unlock()
			return nil
		case msg := <-s.clients[user.Name]:
			err := stream.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (s *ServiceServer) SignUp(ctx context.Context, signUpRequest *messenger.SignUpRequest) (*emptypb.Empty, error) {
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

func (s *ServiceServer) LogIn(ctx context.Context, logInRequest *messenger.LogInRequest) (*messenger.User, error) {
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
	if _, ok := s.clients[dbUser.Name]; ok {
		return nil, convert(errors.ErrUserIsAlreadyOnline)
	}
	return &messenger.User{Name: dbUser.Name}, nil
}
