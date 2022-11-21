package resources

import (
	"context"

	"github.com/dragun-igor/messenger/internal/server/model"
)

type Repository interface {
	InsertMessage(ctx context.Context, message model.Message) error
	CreateUser(ctx context.Context, user model.User) error
	CheckLoginExists(ctx context.Context, user model.User) (bool, error)
	CheckNameExists(ctx context.Context, user model.User) (bool, error)
	LogIn(ctx context.Context, user model.User) (string, string, error)
	Close(ctx context.Context) error
}
