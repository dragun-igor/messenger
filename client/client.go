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

var (
	receiverName = flag.Int64("receiver", 2, "Receiver name for messaging")
	tcpServer    = flag.String("server", ":5400", "Tcp server")
)

func signUp(ctx context.Context, client messengerpb.MessengerServiceClient, scanner *bufio.Scanner) *messengerpb.UserData {
	signUpData := &messengerpb.SignUpData{
		SignInData: &messengerpb.SignInData{},
		UserData:   &messengerpb.UserData{},
	}
	fmt.Print("Your Name: ")
	for signUpData.UserData.Id == 0 {
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
					signUpData.UserData.Name = name
					break
				} else {
					fmt.Println("Name is busy")
				}
			}
		}

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
					signUpData.SignInData.Login = login
					break
				} else {
					fmt.Println("Login is busy")
				}
			}
		}

		for {
			fmt.Print("Your password: ")
			for scanner.Scan() {
				password := scanner.Text()
				if password == "" {
					fmt.Println("Password is empty")
					continue
				}
				signUpData.SignInData.Password = password
				break
			}

			fmt.Print("Your password: ")
			if scanner.Scan() {
				if scanner.Text() != signUpData.SignInData.Password {
					fmt.Println("Passwords do not match")
					continue
				}
				break
			}
		}
		if u, err := client.SignUp(ctx, signUpData); err != nil {
			fmt.Printf("err: %v \n", err)
		} else {
			signUpData.UserData = u
		}
	}
	return signUpData.UserData
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

func receiveMessage(ctx context.Context, client messengerpb.MessengerServiceClient, userData *messengerpb.UserData) {
	stream, err := client.ReceiveMessage(ctx, userData)
	if err != nil {
		log.Fatalf("client.ReceiveMessage(ctx, &userData) throes: %v \n", err)
	}
	for {
		in, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to receive message from user. \nErr: %v \n", err)
		}
		fmt.Printf("%v: %v \n", in.Sender.Name, in.Message)
	}
}

func sendMessage(ctx context.Context, client messengerpb.MessengerServiceClient, message string, user *messengerpb.UserData) {
	stream, err := client.SendMessage(ctx)
	if err != nil {
		log.Printf("Cannot send message: %v \n", err)
	}
	msg := &messengerpb.Message{
		Sender:   user,
		Receiver: &messengerpb.UserData{Id: *receiverName},
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
	userData := &messengerpb.UserData{}
	switch sign {
	case sIn:
		userData = signIn(ctx, client, scanner)
	case sUp:
		userData = signUp(ctx, client, scanner)
	}

	fmt.Printf("Hello, %s!\n", userData.Name)

	go receiveMessage(ctx, client, userData)

	for scanner.Scan() {
		go sendMessage(ctx, client, scanner.Text(), userData)
	}
}
