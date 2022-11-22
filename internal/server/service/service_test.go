package service

import (
	"context"
	"testing"

	"github.com/dragun-igor/messenger/internal/server/model"
	"github.com/dragun-igor/messenger/internal/server/service/mocks"
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

type messageMatcher struct {
	model.Message
}

func (m messageMatcher) Matches(x interface{}) bool {
	m2, ok := x.(model.Message)
	if !ok {
		return false
	}
	return m2 == m.Message
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

func TestMessengerServiceServer(t *testing.T) {
	suite.Run(t, new(MessengerSuiteServer))
}

func (s *MessengerSuiteServer) TestSendMessage() {
	t := s.T()
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
	go func() {
		<-s.service.clients["Receiver"]
	}()
	s.repo.EXPECT().InsertMessage(context.Background(), messageMatcher{modelMessage}.Message).Return(nil)
	resp, err = s.service.SendMessage(context.Background(), protoMessage)
	require.NoError(t, err)
	require.True(t, resp.Sent)
}
