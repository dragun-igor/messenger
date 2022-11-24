package service

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dragun-igor/messenger/internal/pkg/model"
	"github.com/dragun-igor/messenger/proto/messenger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const prefixServiceMessage string = "[SERVICE] "

const (
	nLabel   string = "n"
	noLabel  string = "no"
	yLabel   string = "y"
	yesLabel string = "yes"

	nRusLabel   string = "н"
	noRusLabel  string = "нет"
	yRusLabel   string = "д"
	yesRusLabel string = "да"

	signUpLabel string = "SIGN UP"
	logInLabel  string = "LOG IN"

	loginLabel    string = "Login: "
	nameLabel     string = "Name: "
	passwordLabel string = "Password: "
)

var timeTickerReconnect time.Duration = time.Second * 10

type Service struct {
	client messenger.MessengerServiceClient
	name   string
}

func New(client messenger.MessengerServiceClient) *Service {
	return &Service{client: client}
}

func (c *Service) Serve(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)
	if err := c.auth(ctx, scanner); err != nil {
		return err
	}
	go c.listenMessage(ctx)
	return c.listenScanner(ctx, scanner)
}

func (c *Service) listenScanner(ctx context.Context, scanner *bufio.Scanner) error {
	fmt.Printf(prefixServiceMessage+"Hello, %s!\n", c.name)
	for scanner.Scan() {
		message := scanner.Text()
		switch message {
		case "":
			continue
		default:
			if err := c.sendMessage(ctx, message); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Service) sendMessage(ctx context.Context, message string) error {
	message = strings.TrimSpace(message)
	messageSplit := strings.SplitN(message, " ", 2)
	if len(messageSplit) < 2 {
		fmt.Println(prefixServiceMessage + "Incorrect message. Message should look like \"{username} {message}\"")
		return nil
	}
	response, err := c.client.SendMessage(ctx, &messenger.Message{
		Sender:   c.name,
		Receiver: messageSplit[0],
		Message:  strings.TrimSpace(messageSplit[1]),
	})
	if err != nil {
		return err
	}
	if !response.Sent {
		fmt.Printf(prefixServiceMessage+"User %s is offline!\n", messageSplit[0])
	}
	return nil
}

func (c *Service) listenMessage(ctx context.Context) {
	stream, err := c.client.ReceiveMessage(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	stream.Send(&messenger.User{Name: c.name})

	// reconnecting
	go func() {
		<-stream.Context().Done()
		fmt.Println(prefixServiceMessage + "Connection to server has lost")
		ticker := time.NewTicker(timeTickerReconnect)
		for range ticker.C {
			_, err := c.client.Ping(ctx, &emptypb.Empty{})
			if err != nil {
				fmt.Println(prefixServiceMessage + "Trying to reconnect")
			} else {
				fmt.Println(prefixServiceMessage + "Reconnected")
				go c.listenMessage(ctx)
				break
			}
		}
	}()

	// receiving message from server
	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			return
		}
		fmt.Printf("%v: %v\n", msg.Sender, msg.Message)
	}
}

func (c *Service) auth(ctx context.Context, scanner *bufio.Scanner) error {
BEGIN:
	fmt.Print("Are you already have a account? ")
	if scanner.Scan() {
		text := strings.ToLower(scanner.Text())
		switch text {
		case nLabel, noLabel, nRusLabel, noRusLabel:
			fmt.Println(prefixServiceMessage + signUpLabel)
			if err := c.signUp(ctx, scanner); err != nil {
				return err
			}
			fallthrough
		case yLabel, yesLabel, yRusLabel, yesRusLabel:
			fmt.Println(prefixServiceMessage + logInLabel)
			if err := c.logIn(ctx, scanner); err != nil {
				return err
			}
		default:
			goto BEGIN
		}
	}
	return nil
}

func (c *Service) signUp(ctx context.Context, scanner *bufio.Scanner) error {
BEGIN:
	var authData model.AuthData
	fmt.Print(loginLabel)
	if scanner.Scan() {
		authData.Login = scanner.Text()
	}
	fmt.Print(nameLabel)
	if scanner.Scan() {
		authData.Name = scanner.Text()
	}
	for {
		fmt.Print(passwordLabel)
		if scanner.Scan() {
			authData.Password = scanner.Text()
		}
		fmt.Print(passwordLabel)
		if scanner.Scan() {
			if authData.Password == scanner.Text() {
				break
			}
		}
		fmt.Println(prefixServiceMessage + "Passwords are not matched")
	}
	ve, err := model.Validate(authData)
	if err != nil {
		return err
	}
	if len(ve) > 0 {
		for _, msg := range ve {
			fmt.Println(prefixServiceMessage + msg)
		}
		goto BEGIN
	}
	_, err = c.client.SignUp(ctx, &messenger.SignUpRequest{
		Login:    authData.Login,
		Name:     authData.Name,
		Password: authData.Password,
	})
	if err != nil {
		fmt.Println(stringGRPCError(err))
		goto BEGIN
	}
	return nil
}

func (c *Service) logIn(ctx context.Context, scanner *bufio.Scanner) error {
BEGIN:
	var authData model.AuthData
	fmt.Print(loginLabel)
	if scanner.Scan() {
		authData.Login = scanner.Text()
	}
	fmt.Print(passwordLabel)
	if scanner.Scan() {
		authData.Password = scanner.Text()
	}
	user, err := c.client.LogIn(ctx, &messenger.LogInRequest{
		Login:    authData.Login,
		Password: authData.Password,
	})
	if err != nil {
		fmt.Println(stringGRPCError(err))
		goto BEGIN
	}
	c.name = user.Name
	return nil
}

func stringGRPCError(err error) string {
	grpcErr := status.Convert(err)
	if grpcErr.Code() != codes.Internal {
		return prefixServiceMessage + grpcErr.Proto().Message
	}
	return prefixServiceMessage + grpcErr.Code().String()
}
