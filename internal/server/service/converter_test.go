package service

import (
	"errors"
	"testing"

	errs "github.com/dragun-igor/messenger/pkg/errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConverterSuite struct {
	suite.Suite

	ctrl *gomock.Controller
}

func (s *ConverterSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
}

func (s *ConverterSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestConverterSuite(t *testing.T) {
	suite.Run(t, new(ConverterSuite))
}

func (s *ConverterSuite) TestConvert() {
	cases := []struct {
		err  error
		code codes.Code
		msg  string
	}{
		{
			err: grpcError{
				err:    errors.New("test"),
				status: status.New(codes.Aborted, "message"),
			},
			code: codes.Aborted,
		},
		{
			err:  errs.ErrUserNotFound,
			code: codes.NotFound,
		},
		{
			err:  errs.ErrLoginNameIsBusy,
			code: codes.AlreadyExists,
		},
		{
			err:  errs.ErrUserNameIsBusy,
			code: codes.AlreadyExists,
		},
		{
			err:  errs.ErrIncorrectPassword,
			code: codes.PermissionDenied,
		},
		{
			err:  errors.New("test"),
			code: codes.Internal,
		},
	}

	for n, c := range cases {
		grpcErr := convert(c.err)

		s.Equalf(c.code, grpcErr.GRPCStatus().Code(), "unexpected grpc code in case %d", n)
		s.Equalf(c.err.Error(), grpcErr.Error(), "unexpected error in case %d", n)
	}
}
