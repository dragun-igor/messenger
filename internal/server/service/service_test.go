package service

import (
	"github.com/dragun-igor/messenger/internal/server/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type MessengerSuiteServer struct {
	suite.Suite

	ctrl *gomock.Controller
	repo *mock.MockRepository
}

func (s *MessengerSuiteServer) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.repo = mock.NewMockRepository(s.ctrl)
}

func (s *MessengerSuiteServer) TearDownTest() {
	s.ctrl.Finish()
}

func (s *MessengerSuiteServer) TestInsertMessage_Success() {
	// todo
}

func (s *MessengerSuiteServer) TestInsertMessage_Error() {
	// todo
}

func (s *MessengerSuiteServer) TestCreateUser_Success() {
	// todo
}

func (s *MessengerSuiteServer) TestCreateUser_Error() {
	// todo
}

func (s *MessengerSuiteServer) TestCheckLoginExists_Success() {
	// todo
}

func (s *MessengerSuiteServer) TestCheckLoginExists_Error() {
	// todo
}

func (s *MessengerSuiteServer) TestCheckNameExists_Success() {
	// todo
}

func (s *MessengerSuiteServer) TestCheckNameExists_Error() {
	// todo
}

func (s *MessengerSuiteServer) TestLogIn_Success() {
	// todo
}

func (s *MessengerSuiteServer) TestLogIn_Error() {
	// todo
}

func (s *MessengerSuiteServer) TestClose_Success() {
	// todo
}

func (s *MessengerSuiteServer) TestClose_Error() {
	// todo
}
