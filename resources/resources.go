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

const (
	userTableName    = "users"
	messageTableName = "messages"
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

func (r *Resources) SignIn(ctx context.Context, signInData *messengerpb.SignInData) (int64, string, error) {
	var id int64
	var name string
	query := fmt.Sprintf("SELECT id, first_name FROM %s WHERE login_name = $1 AND pswd = $2", userTableName)
	rows, _ := r.DB.QueryContext(ctx, query, signInData.Login, signInData.Password)
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			return id, name, err
		}
	}
	return id, name, nil
}

func (r *Resources) GetAllMessages(ctx context.Context, user *messengerpb.User) []*messengerpb.Message {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", messageTableName)
	rows, _ := r.DB.QueryContext(ctx, query, user.Id)
	var allMsg []*messengerpb.Message
	for rows.Next() {
		msg := &messengerpb.Message{}
		rows.Scan(msg.Time, msg.Sender, msg.Receiver, msg.Time)
		fmt.Println(msg.Time, msg.Sender, msg.Receiver, msg.Time)
		allMsg = append(allMsg, msg)
	}
	return allMsg

}

func (r *Resources) SendMessage(ctx context.Context, msg *messengerpb.Message) bool {
	query := fmt.Sprintf("INSERT INTO %s (time, sender_id, receiver_id, msg) VALUES ($1, $2, $3, $4);", messageTableName)
	_, err := r.DB.ExecContext(ctx, query, msg.Time.AsTime(), msg.Sender.Id, msg.Receiver.Id, msg.Message)
	if err != nil {
		log.Printf("Query is not executed: %v", err)
		return false
	}
	return true
}

func GetResources(ctx context.Context, config *config.Config) *Resources {
	return &Resources{DB: connectDB(ctx, config)}
}
