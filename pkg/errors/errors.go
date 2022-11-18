package errors

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrUserNameIsBusy = errors.New("user name is busy")
var ErrLoginNameIsBusy = errors.New("login name is busy")
var ErrIncorrectPassword = errors.New("incorrect password")
