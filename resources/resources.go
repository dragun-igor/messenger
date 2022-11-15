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

func (r *Resources) SignUp(ctx context.Context, signUpData *messengerpb.SignUpData) (*messengerpb.UserData, error) {
	var id int64
	var name string
	query := fmt.Sprintf("INSERT INTO %s (login, name, password) VALUES ($1, $2, $3) RETURNING id, name;", userTableName)
	row := r.DB.QueryRowContext(ctx, query, signUpData.SignInData.Login, signUpData.UserData.Name, signUpData.SignInData.Password)
	if err := row.Scan(&id, &name); err != nil {
		return nil, err
	}
	userData := &messengerpb.UserData{
		Id:   id,
		Name: name,
	}
	return userData, nil
}

func (r *Resources) SignIn(ctx context.Context, signInData *messengerpb.SignInData) (*messengerpb.UserData, error) {
	var id int64
	var name string
	query := fmt.Sprintf("SELECT id, name FROM %s WHERE login = $1 AND password = $2", userTableName)
	row := r.DB.QueryRowContext(ctx, query, signInData.Login, signInData.Password)
	if err := row.Scan(&id, &name); err != nil {
		return nil, err
	}
	userData := &messengerpb.UserData{
		Id:   id,
		Name: name,
	}
	return userData, nil
}

func (r *Resources) CheckName(ctx context.Context, checkNameMessage *messengerpb.CheckNameMessage) (*messengerpb.CheckNameAck, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE name = $1", userTableName)
	rows, err := r.DB.QueryContext(ctx, query, checkNameMessage.Name)
	if err != nil {
		return nil, err
	}
	ack := &messengerpb.CheckNameAck{Busy: rows.Next()}
	return ack, nil
}

func (r *Resources) CheckLogin(ctx context.Context, checkLoginMessage *messengerpb.CheckLoginMessage) (*messengerpb.CheckLoginAck, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE login = $1", userTableName)
	rows, err := r.DB.QueryContext(ctx, query, checkLoginMessage.Login)
	if err != nil {
		return nil, err
	}
	ack := &messengerpb.CheckLoginAck{Busy: rows.Next()}
	return ack, nil
}

// func (r *Resources) GetAllMessages(ctx context.Context, user *messengerpb.UserData) []*messengerpb.Message {
// 	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", messageTableName)
// 	rows, _ := r.DB.QueryContext(ctx, query, user.Id)
// 	var allMsg []*messengerpb.Message
// 	for rows.Next() {
// 		msg := &messengerpb.Message{}
// 		rows.Scan(msg.Time, msg.Sender, msg.Receiver, msg.Time)
// 		fmt.Println(msg.Time, msg.Sender, msg.Receiver, msg.Time)
// 		allMsg = append(allMsg, msg)
// 	}
// 	return allMsg

// }

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
