package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/dragun-igor/messenger/messengerpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	sIn uint8 = iota + 1
	sUp
)

var tcpServer = flag.String("server", ":5400", "Tcp server")

var (
	receiver   string
	friends    []string
	requests   []string
	myUserData *messengerpb.User
)

// Ввод и проверка уникальности имени
func signUpName(ctx context.Context, client messengerpb.MessengerServiceClient, scanner *bufio.Scanner) string {
	fmt.Print("Your Name: ")
	for scanner.Scan() {
		name := scanner.Text()
		if name == "" {
			fmt.Println("Name is empty")
			continue
		}
		if ack, err := client.CheckName(ctx, &messengerpb.CheckNameMessage{Name: name}); err != nil {
			fmt.Printf("err: %v \n", err)
		} else {
			fmt.Println(ack.Status)
			if ack.Status == "ok" {
				return name
			}
		}
	}
	return ""
}

// Ввод и проверка уникальности логина
func signUpLogin(ctx context.Context, client messengerpb.MessengerServiceClient, scanner *bufio.Scanner) string {
	fmt.Print("Your login: ")
	for scanner.Scan() {
		login := scanner.Text()
		if login == "" {
			fmt.Println("Name is empty")
			continue
		}
		if ack, err := client.CheckLogin(ctx, &messengerpb.CheckLoginMessage{Login: login}); err != nil {
			fmt.Printf("err: %v \n", err)
		} else {
			fmt.Println(ack.Status)
			if ack.Status == "ok" {
				return login
			}
		}
	}
	return ""
}

// Ввод и проверка пароля
func signUpPassword(scanner *bufio.Scanner) string {
	var password string
	for {
		fmt.Print("Your password: ")
		for scanner.Scan() {
			p := scanner.Text()
			if p == "" {
				fmt.Println("Password is empty")
				continue
			}
			password = p
			break
		}

		fmt.Print("Your password: ")
		if scanner.Scan() {
			if scanner.Text() != password {
				fmt.Println("Passwords do not match")
				continue
			}
			break
		}
	}
	return password
}

// Регистрация
func signUp(ctx context.Context, client messengerpb.MessengerServiceClient, scanner *bufio.Scanner) *messengerpb.User {
	for {
		signUpData := &messengerpb.SignUpData{
			SignInData: &messengerpb.LogInData{
				Login:    signUpLogin(ctx, client, scanner),
				Password: signUpPassword(scanner),
			},
			Name: signUpName(ctx, client, scanner),
		}
		if u, err := client.SignUp(ctx, signUpData); err != nil {
			fmt.Printf("err: %v \n", err)
		} else {
			return u
		}
	}
}

// Авторизация
func logIn(ctx context.Context, client messengerpb.MessengerServiceClient, scanner *bufio.Scanner) *messengerpb.User {
	var user *messengerpb.User
	signInData := &messengerpb.LogInData{}
	for user == nil {
		fmt.Print("Login: ")
		if scanner.Scan() {
			signInData.Login = scanner.Text()
		}
		fmt.Print("Password: ")
		if scanner.Scan() {
			signInData.Password = scanner.Text()
		}
		if u, err := client.LogIn(ctx, signInData); err != nil {
			fmt.Printf("err: %v \nWrong login or password\nTry again\n", err)
		} else {
			user = u
		}
	}
	return user
}

// Получение сообщений
func receiveMessage(ctx context.Context, client messengerpb.MessengerServiceClient) {
	stream, err := client.ReceiveMessage(ctx, myUserData)
	if err != nil {
		log.Fatalf("client.ReceiveMessage(ctx, &userData) throws: %v \n", err)
	}
	for {
		in, err := stream.Recv()
		if err != io.EOF {
			if err != nil {
				log.Fatalf("Failed to receive message from user. \nErr: %v \n", err)
			}
			fmt.Printf("%v: %v \n", in.Sender, in.Message)
		}
	}
}

// Получение архива сообщений с выбранным пользователем
func getMessages(ctx context.Context, client messengerpb.MessengerServiceClient) {
	archive, err := client.GetMessages(ctx, &messengerpb.UsersPair{
		Requester: myUserData.Name,
		Receiver:  receiver,
	})
	if err != nil {
		fmt.Println(err)
	} else {
		for _, msg := range archive.Messages {
			fmt.Printf("%s %s -> %s: %s \n", msg.Time.AsTime().Format("02-01-2006 15:04:05"), msg.Sender, msg.Receiver, msg.Message)
		}
	}
}

// Получение запросов на добавление в друзья
func listenAddToFriendsList(ctx context.Context, client messengerpb.MessengerServiceClient) {
	stream, err := client.ListenAddToFriendsList(ctx, myUserData)
	if err != nil {
		log.Fatalf("client.ListenAddToFriendsList(ctx, &user) throws: %v \n", err)
	}
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			continue
		}
		if err != nil {
			log.Fatalf("Failed to listen requests add to friends list. \nErr: %v \n", err)
		}
		fmt.Printf("User %s request add to friends list \n", in.Name)
	}
}

// Получение подтверждения на запрос
func listenAppendNewFriend(ctx context.Context, client messengerpb.MessengerServiceClient) {
	stream, err := client.ListenAppendNewFriend(ctx, myUserData)
	if err != nil {
		log.Fatalf("client.ListenAddToFriendsList(ctx, &user) throws: %v \n", err)
	}
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			continue
		}
		if err != nil {
			log.Fatalf("Failed to listen append new friend. \nErr: %v \n", err)
		}
		fmt.Printf("User %s added you in friends list \n", in.Name)
		getFriendsList(ctx, client)
	}
}

// Отправление сообщения
func sendMessage(ctx context.Context, client messengerpb.MessengerServiceClient, message string) {
	stream, err := client.SendMessage(ctx)
	if err != nil {
		log.Printf("Cannot send message: %v \n", err)
	}
	msg := &messengerpb.Message{
		Sender:   myUserData.Name,
		Receiver: receiver,
		Message:  message,
	}
	stream.Send(msg)
	ack, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	fmt.Println(ack.Status)
}

// Отправка запроса на добавление в друзья
func requestAddToFriendsList(ctx context.Context, client messengerpb.MessengerServiceClient, name string) {
	ack, err := client.RequestAddToFriendsList(ctx, &messengerpb.UsersPair{Requester: myUserData.Name, Receiver: name})
	if err != nil {
		fmt.Printf("err: %v \n", err)
	} else {
		fmt.Println(ack.Status)
	}
}

// Добавление в друзья
func addToFriendsList(ctx context.Context, client messengerpb.MessengerServiceClient, name string) {
	ack, err := client.AddToFriendsList(ctx, &messengerpb.UsersPair{
		Requester: name,
		Receiver:  myUserData.Name,
	})
	if err != nil {
		fmt.Printf("err: %v", err)
	} else {
		fmt.Println(ack.Status)
	}
}

// Получение списка друзей
func getFriendsList(ctx context.Context, client messengerpb.MessengerServiceClient) {
	friendsList, err := client.GetFriendsList(ctx, myUserData)
	if err != nil {
		log.Println(err)
	}
	friends = friendsList.Friends
}

// Получение списка запросов на добавление в друзья
func getRequestsToFriendsList(ctx context.Context, client messengerpb.MessengerServiceClient) {
	requestsList, err := client.GetRequestsToFriendsList(ctx, myUserData)
	if err != nil {
		log.Println(err)
	}
	requests = requestsList.Requests
}

// Команды пользователя
func command(ctx context.Context, client messengerpb.MessengerServiceClient, message string) {
	str := strings.SplitN(message, " ", 2)
	var command string
	var arg string
	if len(str) > 0 {
		command = str[0]
	}
	if len(str) > 1 {
		arg = str[1]
	}
LOOP:
	switch command {
	case "help", "помощь":
		fmt.Println("HELP GUIDE.")
		fmt.Println("-select or -выбрать - changing message receiver")
		fmt.Println("-archive or -архив - getting message archive (current receiver)")
		fmt.Println("-friends or -друзья - showing friends list")
		fmt.Println("-accept or -принять - add to friends list")
		fmt.Println("-request_to or -запросить - request add to friends list")
		fmt.Println("-requests or -запросы - get requests list")
	case "select", "выбрать":
		for _, name := range friends {
			if name == arg {
				receiver = arg
				fmt.Println("Selected user", receiver)
				break LOOP
			}
		}
		fmt.Printf("User %s not in your friends list \n", arg)
	case "archive", "архив":
		getMessages(ctx, client)
	case "friends", "друзья":
		if len(friends) == 0 {
			fmt.Println("Your friends list is empty")
			break
		}
		fmt.Println("\nYour friends list:")
		for _, name := range friends {
			fmt.Println(name)
		}
	case "accept", "принять":
		addToFriendsList(ctx, client, arg)
		getFriendsList(ctx, client)
	case "request_to", "запросить":
		requestAddToFriendsList(ctx, client, arg)
	case "requests", "запросы":
		getRequestsToFriendsList(ctx, client)
		if len(requests) == 0 {
			fmt.Println("Your requests list is empty")
			break
		}
		fmt.Println("\nYour requests list:")
		for _, name := range requests {
			fmt.Println(name)
		}
	default:
		fmt.Println("Unknown command")
	}
}

func main() {
	flag.Parse()
	fmt.Println("--- CLIENT APP ---")
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(*tcpServer, opts...)
	if err != nil {
		log.Fatalf("Fail to dial: %v \n", err)
	}
	defer conn.Close()
	ctx := context.Background()
	client := messengerpb.NewMessengerServiceClient(conn)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Are you already have a account?")
	var sign uint8
	for sign == 0 {
		if scanner.Scan() {
			str := strings.ToLower(scanner.Text())
			if str == "y" || str == "yes" || str == "д" || str == "да" {
				sign = sIn
			}
			if str == "n" || str == "no" || str == "н" || str == "нет" {
				sign = sUp
			}
		}
	}
	switch sign {
	case sIn:
		myUserData = logIn(ctx, client, scanner)
	case sUp:
		myUserData = signUp(ctx, client, scanner)
	}

	fmt.Printf("Hello, %s!\n", myUserData.Name)
	getFriendsList(ctx, client)
	go receiveMessage(ctx, client)
	go listenAddToFriendsList(ctx, client)
	go listenAppendNewFriend(ctx, client)

	for scanner.Scan() {
		str := scanner.Text()
		switch {
		case str == "":
			continue
		case str[0] == '-':
			go command(ctx, client, str[1:])
		default:
			if receiver != "" {
				go sendMessage(ctx, client, scanner.Text())
			} else {
				fmt.Println("Select message receiver")
			}
		}
	}
}
