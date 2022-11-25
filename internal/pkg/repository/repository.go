package repository

import (
	"context"
	"fmt"

	"github.com/dragun-igor/messenger/internal/pkg/model"
)

const (
	usersTable    string = "users"
	messagesTable string = "messages"
)

type Repository struct {
	db Postgres
}

func New(db Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) InsertMessage(ctx context.Context, message model.Message) error {
	query := fmt.Sprintf("INSERT INTO %s VALUES ($1, $2, $3);", messagesTable)
	_, err := r.db.Exec(ctx, query, message.Sender, message.Receiver, message.Message)
	return err
}

func (r *Repository) CreateUser(ctx context.Context, user model.AuthData) error {
	query := fmt.Sprintf("INSERT INTO %s VALUES ($1, $2, $3);", usersTable)
	_, err := r.db.Exec(ctx, query, user.Login, user.Name, user.Password)
	return err
}

func (r *Repository) CheckLoginExists(ctx context.Context, user model.AuthData) (bool, error) {
	var ok bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT * FROM %s WHERE login = $1)", usersTable)
	row := r.db.QueryRow(ctx, query, user.Login)
	err := row.Scan(&ok)
	return !ok, err
}

func (r *Repository) CheckNameExists(ctx context.Context, user model.AuthData) (bool, error) {
	var ok bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT * FROM %s WHERE name = $1)", usersTable)
	row := r.db.QueryRow(ctx, query, user.Name)
	err := row.Scan(&ok)
	return !ok, err
}

func (r *Repository) GetUser(ctx context.Context, user model.AuthData) (model.AuthData, error) {
	var login string
	var name string
	var password string
	query := fmt.Sprintf("SELECT * FROM %s WHERE login = $1", usersTable)
	row := r.db.QueryRow(ctx, query, user.Login)
	err := row.Scan(&login, &name, &password)
	return model.AuthData{Login: login, Name: name, Password: password}, err
}
