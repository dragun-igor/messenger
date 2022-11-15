package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
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
	receiver   string = "Ira"
	myUserData *messengerpb.UserData
)

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
			if !ack.Busy {
				return name
			} else {
				fmt.Println("Name is busy")
			}
		}
	}
	return ""
}

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
			if !ack.Busy {
				return login
			} else {
				fmt.Println("Login is busy")
			}
		}
	}
	return ""
}

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

func signUp(ctx context.Context, client messengerpb.MessengerServiceClient, scanner *bufio.Scanner) *messengerpb.UserData {
	for {
		signUpData := &messengerpb.SignUpData{
			SignInData: &messengerpb.SignInData{
				Login:    signUpLogin(ctx, client, scanner),
				Password: signUpPassword(scanner),
			},
			UserData: &messengerpb.UserData{
				Name: signUpName(ctx, client, scanner),
			},
		}
		if u, err := client.SignUp(ctx, signUpData); err != nil {
			fmt.Printf("err: %v \n", err)
		} else {
			return u
		}
	}
}

func signIn(ctx context.Context, client messengerpb.MessengerServiceClient, scanner *bufio.Scanner) *messengerpb.UserData {
	var userData *messengerpb.UserData
	signInData := &messengerpb.SignInData{}
	for userData == nil {
		fmt.Print("Login: ")
		if scanner.Scan() {
			signInData.Login = scanner.Text()
		}
		fmt.Print("Password: ")
		if scanner.Scan() {
			signInData.Password = scanner.Text()
		}
		if u, err := client.SignIn(ctx, signInData); err != nil {
			fmt.Printf("err: %v \nWrong login or password\nTry again\n", err)
		} else {
			userData = u
		}
	}
	return userData
}

func receiveMessage(ctx context.Context, client messengerpb.MessengerServiceClient) {
	stream, err := client.ReceiveMessage(ctx, myUserData)
	if err != nil {
		log.Fatalf("client.ReceiveMessage(ctx, &userData) throws: %v \n", err)
	}
	for {
		in, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to receive message from user. \nErr: %v \n", err)
		}
		fmt.Printf("%v: %v \n", in.Sender.Name, in.Message)
	}
}

func listenAddToFriendsList(ctx context.Context, client messengerpb.MessengerServiceClient) {
	stream, err := client.ListenAddToFriendsList(ctx, myUserData)
	if err != nil {
		log.Fatalf("client.ListenAddToFriendsList(ctx, &userData) throws: %v \n", err)
	}
	for {
		in, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to listen requests add to friends list. \nErr: %v \n", err)
		}
		fmt.Printf("User %s request add to friends list \n", in.Name)
	}
}

func sendMessage(ctx context.Context, client messengerpb.MessengerServiceClient, message string, user *messengerpb.UserData) {
	stream, err := client.SendMessage(ctx)
	if err != nil {
		log.Printf("Cannot send message: %v \n", err)
	}
	msg := &messengerpb.Message{
		Sender:   user,
		Receiver: &messengerpb.UserData{Name: receiver},
		Message:  message,
	}
	stream.Send(msg)
	ack, _ := stream.CloseAndRecv()
	fmt.Printf("Message status: %v \n", ack)
}

func requestAddToFriendsList(ctx context.Context, client messengerpb.MessengerServiceClient, name string) {
	_, err := client.RequestAddToFriendsList(ctx, &messengerpb.RequestAddToFriendsListMessage{Requester: myUserData, Receiver: &messengerpb.UserData{Name: name}})
	if err != nil {
		fmt.Printf("err: %v \n", err)
	} else {
		fmt.Printf("Request sent")
	}
}

func command(ctx context.Context, client messengerpb.MessengerServiceClient, message string) {
	switch {
	case message == "помощь", message == "help":
		fmt.Println("HELP GUIDE.")
		fmt.Println("-receiver or -получатель - changing message receiver")
		fmt.Println("-archive or -архив - getting message archive (current receiver)")
	case strings.SplitN(message, " ", 2)[0] == "receiver", strings.Split(message, " ")[0] == "получатель":
		fmt.Printf("New reciver %s \n", strings.SplitN(message, " ", 2)[1])
		receiver = strings.SplitN(message, " ", 2)[1]
	case message == "friends", message == "друзья":
	case strings.SplitN(message, " ", 2)[0] == "request":
		requestAddToFriendsList(ctx, client, strings.SplitN(message, " ", 2)[1])
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
		myUserData = signIn(ctx, client, scanner)
	case sUp:
		myUserData = signUp(ctx, client, scanner)
	}

	fmt.Printf("Hello, %s!\n", myUserData.Name)

	go receiveMessage(ctx, client)
	go listenAddToFriendsList(ctx, client)

	for scanner.Scan() {
		str := scanner.Text()
		if str[0] == '-' {
			go command(ctx, client, str[1:])
		} else {
			go sendMessage(ctx, client, scanner.Text(), myUserData)
		}
	}
}
