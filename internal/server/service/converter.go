package service

import (
	"github.com/dragun-igor/messenger/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCError interface {
	Error() string
	GRPCStatus() *status.Status
	Unwrap() error
}

type grpcError struct {
	err    error
	status *status.Status
}

func convert(err error) GRPCError {
	if v, ok := err.(grpcError); ok {
		return v
	}
	switch err {
	case errors.ErrUserNotFound:
		return newGRPCError(err, codes.NotFound)
	case errors.ErrUserNameIsBusy, errors.ErrLoginNameIsBusy:
		return newGRPCError(err, codes.AlreadyExists)
	case errors.ErrIncorrectPassword:
		return newGRPCError(err, codes.PermissionDenied)
	default:
		return newGRPCError(err, codes.Internal)
	}
}

func newGRPCError(err error, code codes.Code) grpcError {
	return grpcError{
		err:    err,
		status: status.New(code, err.Error()),
	}
}

func (e grpcError) Error() string {
	return e.err.Error()
}

func (e grpcError) GRPCStatus() *status.Status {
	return e.status
}

func (e grpcError) Unwrap() error {
	return e.err
}