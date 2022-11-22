package service

import (
	"github.com/dragun-igor/messenger/internal/server/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type MessengerSuiteServer struct {
	suite.Suite

	ctrl *gomock.Controller
	repo *mocks.MockRepository
}

func (s *MessengerSuiteServer) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.repo = mocks.NewMockRepository(s.ctrl)
}

func (s *MessengerSuiteServer) TearDownTest() {
	s.ctrl.Finish()
}

func (s *MessengerSuiteServer) TestSignUp_Success() {

}

func (s *MessengerSuiteServer) TestSignUp_Error() {
	// todo
}

func (s *MessengerSuiteServer) TestLogIn_Success() {
	// todo
}

func (s *MessengerSuiteServer) TestLogIn_Error() {
	// todo
}

func (s *MessengerSuiteServer) TestSendMessage_Success() {
	// todo
}

func (s *MessengerSuiteServer) TestSendMessage_Error() {
	// todo
}
