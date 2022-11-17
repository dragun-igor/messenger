package resources

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/dragun-igor/messenger/config"
	"github.com/dragun-igor/messenger/messengerpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	_ "github.com/lib/pq"
)

// Названия таблиц
const (
	usersTableName            = "users"
	messagesTableName         = "messages"
	friendsListTableName      = "friends"
	requestToFriendsTableName = "requests_to_friends_list"
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
		log.Println("connection to database has closed")
	}()
	return db
}

func (r *Resources) SignUp(ctx context.Context, signUpData *messengerpb.SignUpData) (*messengerpb.User, error) {
	var name string
	query := fmt.Sprintf("INSERT INTO %s (login, name, password) VALUES ($1, $2, $3) RETURNING name;", usersTableName)
	row := r.DB.QueryRowContext(ctx, query, signUpData.SignInData.Login, signUpData.Name, signUpData.SignInData.Password)
	if err := row.Scan(&name); err != nil {
		return nil, err
	}
	userData := &messengerpb.User{
		Name: name,
	}
	return userData, nil
}

func (r *Resources) SignIn(ctx context.Context, logInData *messengerpb.LogInData) (*messengerpb.User, error) {
	var name string
	query := fmt.Sprintf("SELECT name FROM %s WHERE login = $1 AND password = $2;", usersTableName)
	row := r.DB.QueryRowContext(ctx, query, logInData.Login, logInData.Password)
	if err := row.Scan(&name); err != nil {
		return nil, err
	}
	userData := &messengerpb.User{
		Name: name,
	}
	return userData, nil
}

func (r *Resources) CheckName(ctx context.Context, checkNameMessage *messengerpb.CheckNameMessage) (*messengerpb.MessageAck, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE name = $1;", usersTableName)
	row := r.DB.QueryRowContext(ctx, query, checkNameMessage.Name)
	var number int
	if err := row.Scan(&number); err != nil {
		return nil, err
	}
	ack := &messengerpb.MessageAck{}
	if number > 0 {
		ack.Status = "Name is busy"
	} else {
		ack.Status = "Ok"
	}
	return ack, nil
}

func (r *Resources) CheckLogin(ctx context.Context, checkLoginMessage *messengerpb.CheckLoginMessage) (*messengerpb.MessageAck, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE login = $1;", usersTableName)
	row := r.DB.QueryRowContext(ctx, query, checkLoginMessage.Login)
	var number int
	if err := row.Scan(&number); err != nil {
		return nil, err
	}
	ack := &messengerpb.MessageAck{}
	if number > 0 {
		ack.Status = "Login is busy"
	} else {
		ack.Status = "Ok"
	}
	return ack, nil
}

func (r *Resources) RequestAddToFriendsList(ctx context.Context, usersPair *messengerpb.UsersPair) (*messengerpb.MessageAck, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE requester = $1 AND receiver = $2;", requestToFriendsTableName)
	row := r.DB.QueryRowContext(ctx, query, usersPair.Requester, usersPair.Receiver)
	ack := &messengerpb.MessageAck{}
	var number int
	if err := row.Scan(&number); err != nil {
		return nil, err
	}
	if number > 0 {
		ack.Status = "Request is active already"
		return ack, nil
	}
	query = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user1 = $1 AND user2 = $2;", friendsListTableName)
	row = r.DB.QueryRowContext(ctx, query, usersPair.Requester, usersPair.Receiver)
	if err := row.Scan(&number); err != nil {
		return nil, err
	}
	if number > 0 {
		ack.Status = "User in your friends list already"
		return ack, nil
	}
	query = fmt.Sprintf("INSERT INTO %s (requester, receiver) VALUES ($1, $2);", requestToFriendsTableName)
	_, err := r.DB.ExecContext(ctx, query, usersPair.Requester, usersPair.Receiver)
	if err != nil {
		return nil, err
	}
	ack.Status = "Request sent"
	return ack, nil
}

func (r *Resources) AddToFriendsList(ctx context.Context, usersPair *messengerpb.UsersPair) (*messengerpb.MessageAck, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE name = $1", usersTableName)
	row := r.DB.QueryRowContext(ctx, query, usersPair.Requester)
	var number int
	ack := &messengerpb.MessageAck{}
	if err := row.Scan(&number); err != nil {
		return nil, err
	}
	if number == 0 {
		ack.Status = "User not found"
		return ack, nil
	}
	var res int
	query = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE requester = $1 and receiver = $2", requestToFriendsTableName)
	row = r.DB.QueryRowContext(ctx, query, usersPair.Requester, usersPair.Receiver)
	if err := row.Scan(&res); err != nil {
		return nil, err
	}
	if res == 0 {
		ack.Status = "User doesn't request add to friends list"
		return ack, nil
	}
	query = fmt.Sprintf("INSERT INTO %s VALUES ($1, $2), ($2, $1);", friendsListTableName)
	_, err := r.DB.ExecContext(ctx, query, usersPair.Requester, usersPair.Receiver)
	if err != nil {
		return nil, err
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE requster = $1 AND receiver = $2", requestToFriendsTableName)
	_, err = r.DB.ExecContext(ctx, query, usersPair.Requester, usersPair.Receiver)
	if err != nil {
		log.Println(err)
	}
	ack.Status = "Ok"
	return ack, nil
}

func (r *Resources) GetAllUsers(ctx context.Context) ([]string, error) {
	var usersNumber int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", usersTableName)
	row := r.DB.QueryRowContext(ctx, query)
	if err := row.Scan(&usersNumber); err != nil {
		return nil, err
	}
	users := make([]string, 0, usersNumber)
	query = fmt.Sprintf("SELECT name FROM %s", usersTableName)
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		users = append(users, name)
	}
	return users, nil
}

func (r *Resources) GetFriendsList(ctx context.Context, user *messengerpb.User) (*messengerpb.FriendsList, error) {
	query := fmt.Sprintf("SELECT COUNT(user2) FROM %s WHERE user1 = $1;", friendsListTableName)
	row := r.DB.QueryRowContext(ctx, query, user.Name)
	var number int
	if err := row.Scan(&number); err != nil {
		return nil, err
	}
	if number == 0 {
		return &messengerpb.FriendsList{}, nil
	}
	friends := make([]string, 0, number)
	query = fmt.Sprintf("SELECT user1 FROM %s WHERE user2 = $1;", friendsListTableName)
	rows, err := r.DB.QueryContext(ctx, query, user.Name)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var name string
		rows.Scan(&name)
		friends = append(friends, name)
	}
	return &messengerpb.FriendsList{Friends: friends}, nil
}

func (r *Resources) GetRequestsToFriendsList(ctx context.Context, user *messengerpb.User) (*messengerpb.Requests, error) {
	query := fmt.Sprintf("SELECT COUNT(requester) FROM %s WHERE receiver = $1;", requestToFriendsTableName)
	row := r.DB.QueryRowContext(ctx, query, user.Name)
	var number int
	if err := row.Scan(&number); err != nil {
		return nil, err
	}
	if number == 0 {
		return &messengerpb.Requests{}, nil
	}
	requests := make([]string, 0, number)
	query = fmt.Sprintf("SELECT requester FROM %s WHERE receiver = $1;", requestToFriendsTableName)
	rows, err := r.DB.QueryContext(ctx, query, user.Name)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var name string
		rows.Scan(&name)
		requests = append(requests, name)
	}
	return &messengerpb.Requests{Requests: requests}, nil
}

func (r *Resources) GetMessages(ctx context.Context, usersPair *messengerpb.UsersPair) (*messengerpb.MessageArchive, error) {
	query := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE (sender = $1 AND receiver = $2) OR (sender = $2 AND receiver = $1);", messagesTableName)
	row := r.DB.QueryRowContext(ctx, query, usersPair.Requester, usersPair.Receiver)
	var number int
	if err := row.Scan(&number); err != nil {
		return nil, err
	}
	if number == 0 {
		return &messengerpb.MessageArchive{Messages: []*messengerpb.Message{}}, nil
	}
	messages := make([]*messengerpb.Message, 0, number)
	query = fmt.Sprintf("SELECT time, sender, receiver, msg FROM %s WHERE (sender = $1 AND receiver = $2) OR (sender = $2 AND receiver = $1)", messagesTableName)
	rows, err := r.DB.QueryContext(ctx, query, usersPair.Requester, usersPair.Receiver)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var t time.Time
		var senderName string
		var receiverName string
		var text string
		rows.Scan(&t, &senderName, &receiverName, &text)
		message := &messengerpb.Message{
			Time:     timestamppb.New(t),
			Sender:   senderName,
			Receiver: receiverName,
			Message:  text,
		}
		messages = append(messages, message)
	}
	return &messengerpb.MessageArchive{Messages: messages}, nil
}

func (r *Resources) SendMessage(ctx context.Context, msg *messengerpb.Message) error {
	query := fmt.Sprintf("INSERT INTO %s (time, sender, receiver, msg) VALUES ($1, $2, $3, $4);", messagesTableName)
	_, err := r.DB.ExecContext(ctx, query, msg.Time.AsTime(), msg.Sender, msg.Receiver, msg.Message)
	return err
}

func GetResources(ctx context.Context, config *config.Config) *Resources {
	return &Resources{DB: connectDB(ctx, config)}
}
