package service

import (
	"context"

	"github.com/dragun-igor/messenger/internal/server/model"
)

type Repository interface {
	InsertMessage(ctx context.Context, message model.Message) error
	CreateUser(ctx context.Context, user model.AuthData) error
	CheckLoginExists(ctx context.Context, user model.AuthData) (bool, error)
	CheckNameExists(ctx context.Context, user model.AuthData) (bool, error)
	GetUser(ctx context.Context, user model.AuthData) (model.AuthData, error)
}
