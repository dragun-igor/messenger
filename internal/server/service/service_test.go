package service

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/dragun-igor/messenger/internal/pkg/model"
	"github.com/dragun-igor/messenger/internal/server/service/mocks"
	"github.com/dragun-igor/messenger/pkg/errors"
	"github.com/dragun-igor/messenger/proto/messenger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type MessengerSuiteServer struct {
	suite.Suite

	ctrl    *gomock.Controller
	repo    *mocks.MockRepository
	service *MessengerServiceServer
}

type authMatcher struct {
	model.AuthData
}

func (a authMatcher) Matches(x interface{}) bool {
	a2, ok := x.(model.AuthData)
	if !ok {
		return false
	}
	return a2.Login == a.AuthData.Login && a2.Name == a.AuthData.Name && a2.IsPasswordCorrect(a.AuthData.Password)
}

func (a authMatcher) String() string {
	return ""
}

func (s *MessengerSuiteServer) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.repo = mocks.NewMockRepository(s.ctrl)
	grpc := grpc.NewServer([]grpc.ServerOption{}...)
	service := NewMessengerServiceServer(context.Background(), s.repo)
	messenger.RegisterMessengerServiceServer(grpc, service)
	s.service = service
	s.service.clients["Receiver"] = make(chan *messenger.Message)
}

func (s *MessengerSuiteServer) TearDownTest() {
	s.ctrl.Finish()
}

func TestMessengerSuiteServer(t *testing.T) {
	suite.Run(t, new(MessengerSuiteServer))
}

func (s *MessengerSuiteServer) TestSendMessage() {
	t := s.T()
	ctx := context.Background()
	protoMessage := &messenger.Message{
		Sender:   "Sender",
		Receiver: "UnknownReceiver",
		Message:  "Message",
	}
	modelMessage := model.Message{
		Sender:   "Sender",
		Receiver: "Receiver",
		Message:  "Message",
	}
	resp, err := s.service.SendMessage(context.Background(), protoMessage)
	require.False(t, resp.Sent)
	require.NoError(t, err)
	protoMessage.Receiver = "Receiver"
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		msg := <-s.service.clients["Receiver"]
		require.Equal(t, protoMessage.Sender, msg.Sender)
		require.Equal(t, protoMessage.Receiver, msg.Receiver)
		require.Equal(t, protoMessage.Message, msg.Message)
	}()
	s.repo.EXPECT().InsertMessage(ctx, modelMessage).Return(nil)
	resp, err = s.service.SendMessage(ctx, protoMessage)
	require.NoError(t, err)
	require.True(t, resp.Sent)
	wg.Wait()
}

func (s *MessengerSuiteServer) TestAuth() {
	t := s.T()
	ctx := context.Background()
	protoAuth := &messenger.SignUpRequest{
		Login:    "Login",
		Name:     "Name",
		Password: "Password",
	}
	modelAuth := model.AuthData{
		Login: "Login",
		Name:  "Name",
	}
	s.repo.EXPECT().CheckLoginExists(ctx, modelAuth).Return(false, nil)
	_, err := s.service.SignUp(ctx, protoAuth)
	require.Error(t, convert(errors.ErrLoginNameIsBusy), err)

	testErr := fmt.Errorf("test error")
	s.repo.EXPECT().CheckLoginExists(ctx, modelAuth).Return(false, testErr)
	_, err = s.service.SignUp(ctx, protoAuth)
	require.Error(t, convert(testErr), err)

	s.repo.EXPECT().CheckLoginExists(ctx, modelAuth).Return(true, nil)
	s.repo.EXPECT().CheckNameExists(ctx, modelAuth).Return(false, nil)
	_, err = s.service.SignUp(ctx, protoAuth)
	require.Error(t, convert(errors.ErrUserNameIsBusy), err)

	s.repo.EXPECT().CheckLoginExists(ctx, modelAuth).Return(true, nil)
	s.repo.EXPECT().CheckNameExists(ctx, modelAuth).Return(false, testErr)
	_, err = s.service.SignUp(ctx, protoAuth)
	require.Error(t, convert(testErr), err)

	s.repo.EXPECT().CheckLoginExists(ctx, modelAuth).Return(true, nil)
	s.repo.EXPECT().CheckNameExists(ctx, modelAuth).Return(true, nil)
	modelAuth.Password = "Password"
	s.repo.EXPECT().CreateUser(ctx, authMatcher{modelAuth}).Return(nil)
	_, err = s.service.SignUp(ctx, protoAuth)
	require.NoError(t, err)
}
