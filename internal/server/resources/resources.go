package resources

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/dragun-igor/messenger/config"
	"github.com/dragun-igor/messenger/internal/server/model"
	"github.com/jackc/pgx/v5"
)

const (
	usersTable    string = "users"
	messagesTable string = "messages"
)

type PostgresDB struct {
	*pgx.Conn
}

func InitPostgresDB(ctx context.Context, migrationsPath string, config *config.Config) (PostgresDB, error) {
	db, err := connectDB(ctx, migrationsPath, config)
	return PostgresDB{db}, err
}

func connectDB(ctx context.Context, migrationsPath string, config *config.Config) (*pgx.Conn, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	db, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	log.Println("connection to db has opened")

	query, err := getMigrationQuery(migrationsPath)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(ctx, query); err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		_ = db.Close(ctx)
		log.Println("connection to db has closed")
	}()
	return db, nil
}

func getMigrationQuery(migrationsPath string) (string, error) {
	file, err := os.Open(migrationsPath)
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (db PostgresDB) InsertMessage(ctx context.Context, message model.Message) error {
	query := fmt.Sprintf("INSERT INTO %s VALUES ($1, $2, $3);", messagesTable)
	_, err := db.Exec(ctx, query, message.Sender, message.Receiver, message.Message)
	return err
}

func (db PostgresDB) CreateUser(ctx context.Context, user model.User) error {
	query := fmt.Sprintf("INSERT INTO %s VALUES ($1, $2, $3);", usersTable)
	_, err := db.Exec(ctx, query, user.Login, user.Name, user.Password)
	return err
}

func (db PostgresDB) CheckLoginExists(ctx context.Context, user model.User) (bool, error) {
	var ok bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT * FROM %s WHERE login = $1)", usersTable)
	row := db.QueryRow(ctx, query, user.Login)
	err := row.Scan(&ok)
	return !ok, err
}

func (db PostgresDB) CheckNameExists(ctx context.Context, user model.User) (bool, error) {
	var ok bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT * FROM %s WHERE name = $1)", usersTable)
	row := db.QueryRow(ctx, query, user.Name)
	err := row.Scan(&ok)
	return !ok, err
}

func (db PostgresDB) LogIn(ctx context.Context, user model.User) (string, string, error) {
	var name string
	var password string
	query := fmt.Sprintf("SELECT name, password FROM %s WHERE login = $1", usersTable)
	row := db.QueryRow(ctx, query, user.Login)
	err := row.Scan(&name, &password)
	return name, password, err
}
