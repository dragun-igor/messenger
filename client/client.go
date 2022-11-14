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
	senderName   = flag.Int64("sender", 1, "Sender name for messaging")
	receiverName = flag.Int64("receiver", 2, "Receiver name for messaging")
	tcpServer    = flag.String("server", ":5400", "Tcp server")
)

func signIn(ctx context.Context, client messengerpb.MessengerServiceClient, signInData *messengerpb.SignInData) *messengerpb.UserID {
	userId, err := client.SignIn(ctx, signInData)
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	return userId
}

func receiveMessage(ctx context.Context, client messengerpb.MessengerServiceClient) {
	userID := &messengerpb.UserID{Id: *senderName}
	stream, err := client.ReceiveMessage(ctx, userID)
	if err != nil {
		log.Fatalf("client.ReceiveMessage(ctx, &userID) throes: %v \n", err)
	}
	for {
		in, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to receive message from user. \nErr: %v \n", err)
		}
		fmt.Printf("%v: %v \n", in.Sender, in.Message)
	}
}

func sendMessage(ctx context.Context, client messengerpb.MessengerServiceClient, message string) {
	stream, err := client.SendMessage(ctx)
	if err != nil {
		log.Printf("Cannot send message: %v \n", err)
	}
	msg := &messengerpb.Message{
		Sender:   &messengerpb.UserID{Id: *senderName},
		Receiver: &messengerpb.UserID{Id: *receiverName},
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

	userId := &messengerpb.UserID{}
	signInData := &messengerpb.SignInData{}
	scanner := bufio.NewScanner(os.Stdin)
	for userId.Id == 0 {
		fmt.Println("Login: ")
		if scanner.Scan() {
			signInData.Login = scanner.Text()
		}
		fmt.Println("Password: ")
		if scanner.Scan() {
			signInData.Password = scanner.Text()
		}
		userId = signIn(ctx, client, signInData)
	}

	go receiveMessage(ctx, client)

	for scanner.Scan() {
		go sendMessage(ctx, client, scanner.Text())
	}
}
