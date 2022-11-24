package errors

import "errors"

var (
	ErrUserNameIsBusy      = errors.New("user name is busy")
	ErrLoginNameIsBusy     = errors.New("login name is busy")
	ErrIncorrectPassword   = errors.New("incorrect password")
	ErrUserIsAlreadyOnline = errors.New("user is already online")
)
