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
	"github.com/jackc/pgx/v5"
)

type Connection struct {
	*pgx.Conn
}

func NewConnection(ctx context.Context, config *config.Config) (Connection, error) {
	db, err := connectDB(ctx, config)
	return Connection{db}, err
}

func connectDB(ctx context.Context, config *config.Config) (*pgx.Conn, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)
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
