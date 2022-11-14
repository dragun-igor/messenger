package resources

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/dragun-igor/messenger/config"
	"github.com/dragun-igor/messenger/messengerpb"

	_ "github.com/lib/pq"
)

type Resources struct {
	DB *sql.DB
}

func connectDB(ctx context.Context, config *config.Config) *sql.DB {
	fmt.Println(*config)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("cannot connect to database: %v \n", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("database doesn't pinged: %v \n", err)
	}
	log.Println("succesfully connected to database")
	go func() {
		<-ctx.Done()
		_ = db.Close()
		log.Println("Connection to database has closed")
	}()
	return db
}

func (r *Resources) SendMessage(msg *messengerpb.Message) {
	res, err := r.DB.Exec("INSERT INTO messages (sender_id, receiver_id, msg) VALUES ($1, $2, $3);", msg.Sender.Id, msg.Receiver.Id, msg.Message)
	if err != nil {
		log.Printf("Query is not executed: %v", err)
	} else {
		log.Println(res)
	}
}

func GetResources(ctx context.Context, config *config.Config) *Resources {
	return &Resources{DB: connectDB(ctx, config)}
}
