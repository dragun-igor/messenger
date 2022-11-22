package resources

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/dragun-igor/messenger/config"
	"github.com/dragun-igor/messenger/internal/pkg/model"
	"github.com/jackc/pgx/v5"
)

const (
	usersTable    string = "users"
	messagesTable string = "messages"
)

type PostgresDB struct {
	*pgx.Conn
}

func InitPostgresDB(ctx context.Context, config *config.Config) (PostgresDB, error) {
	db, err := connectDB(ctx, config)
	return PostgresDB{db}, err
}

func connectDB(ctx context.Context, config *config.Config) (*pgx.Conn, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	db, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	log.Println("connection to db has opened")

	queries, err := getMigrationsQuery(config.MigrationsPath)
	if err != nil {
		return nil, err
	}
	for _, query := range queries {
		if _, err := db.Exec(ctx, query); err != nil {
			return nil, err
		}
	}
	return db, nil
}

func getMigrationsQuery(migrationsPath string) ([]string, error) {
	migrationsQuery := []string{}
	dirEntries, err := os.ReadDir(migrationsPath)
	if err != nil {
		return nil, err
	}
	migrationFileNames := []string{}
	for _, de := range dirEntries {
		if !de.IsDir() && strings.Contains(de.Name(), ".up.sql") {
			migrationFileNames = append(migrationFileNames, de.Name())
		}
	}
	sort.Slice(migrationFileNames, func(i, j int) bool {
		return migrationFileNames[i] < migrationFileNames[j]
	})
	for _, fileName := range migrationFileNames {
		file, err := os.Open(migrationsPath + "/" + fileName)
		if err != nil {
			return nil, err
		}
		b, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		migrationsQuery = append(migrationsQuery, string(b))
	}

	return migrationsQuery, nil
}

func (db PostgresDB) InsertMessage(ctx context.Context, message model.Message) error {
	query := fmt.Sprintf("INSERT INTO %s VALUES ($1, $2, $3);", messagesTable)
	_, err := db.Exec(ctx, query, message.Sender, message.Receiver, message.Message)
	return err
}

func (db PostgresDB) CreateUser(ctx context.Context, user model.AuthData) error {
	query := fmt.Sprintf("INSERT INTO %s VALUES ($1, $2, $3);", usersTable)
	_, err := db.Exec(ctx, query, user.Login, user.Name, user.Password)
	return err
}

func (db PostgresDB) CheckLoginExists(ctx context.Context, user model.AuthData) (bool, error) {
	var ok bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT * FROM %s WHERE login = $1)", usersTable)
	row := db.QueryRow(ctx, query, user.Login)
	err := row.Scan(&ok)
	return !ok, err
}

func (db PostgresDB) CheckNameExists(ctx context.Context, user model.AuthData) (bool, error) {
	var ok bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT * FROM %s WHERE name = $1)", usersTable)
	row := db.QueryRow(ctx, query, user.Name)
	err := row.Scan(&ok)
	return !ok, err
}

func (db PostgresDB) GetUser(ctx context.Context, user model.AuthData) (model.AuthData, error) {
	var login string
	var name string
	var password string
	query := fmt.Sprintf("SELECT * FROM %s WHERE login = $1", usersTable)
	row := db.QueryRow(ctx, query, user.Login)
	err := row.Scan(&login, &name, &password)
	return model.AuthData{Login: login, Name: name, Password: password}, err
}
