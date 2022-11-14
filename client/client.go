package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dragun-igor/messenger/messengerpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	receiverName = flag.Int64("receiver", 2, "Receiver name for messaging")
	tcpServer    = flag.String("server", ":5400", "Tcp server")
)

func signIn(ctx context.Context, client messengerpb.MessengerServiceClient, signInData *messengerpb.SignInData) *messengerpb.User {
	user, err := client.SignIn(ctx, signInData)
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	return user
}

func receiveMessage(ctx context.Context, client messengerpb.MessengerServiceClient, user *messengerpb.User) {
	stream, err := client.ReceiveMessage(ctx, user)
	if err != nil {
		log.Fatalf("client.ReceiveMessage(ctx, &user) throes: %v \n", err)
	}
	for {
		in, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to receive message from user. \nErr: %v \n", err)
		}
		fmt.Printf("%v: %v \n", in.Sender.FirstName, in.Message)
	}
}

func sendMessage(ctx context.Context, client messengerpb.MessengerServiceClient, message string, user *messengerpb.User) {
	stream, err := client.SendMessage(ctx)
	if err != nil {
		log.Printf("Cannot send message: %v \n", err)
	}
	msg := &messengerpb.Message{
		Sender:   user,
		Receiver: &messengerpb.User{Id: *receiverName},
		Message:  message,
	}
	stream.Send(msg)
	ack, _ := stream.CloseAndRecv()
	fmt.Printf("Message status: %v \n", ack)
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

	user := &messengerpb.User{}
	signInData := &messengerpb.SignInData{}
	scanner := bufio.NewScanner(os.Stdin)
	for user.Id == 0 {
		fmt.Println("Login: ")
		if scanner.Scan() {
			signInData.Login = scanner.Text()
		}
		fmt.Println("Password: ")
		if scanner.Scan() {
			signInData.Password = scanner.Text()
		}
		user = signIn(ctx, client, signInData)
	}
	fmt.Printf("Hello, %s!\n", user.FirstName)

	go receiveMessage(ctx, client, user)

	for scanner.Scan() {
		go sendMessage(ctx, client, scanner.Text(), user)
	}
}
