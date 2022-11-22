package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/dragun-igor/messenger/internal/pkg/model"
	"github.com/dragun-igor/messenger/proto/messenger"
)

const prefixServiceMessage string = "[SERVICE] "

type MessengerServiceClient struct {
	client messenger.MessengerServiceClient
	name   string
}

func NewClientService(client messenger.MessengerServiceClient) *MessengerServiceClient {
	return &MessengerServiceClient{client: client}
}

func (c *MessengerServiceClient) Serve(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)
	if err := c.auth(ctx, scanner); err != nil {
		return err
	}
	if err := c.connect(ctx); err != nil {
		return err
	}
	return c.listenScanner(ctx, scanner)
}

func (c *MessengerServiceClient) listenScanner(ctx context.Context, scanner *bufio.Scanner) error {
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

func (c *MessengerServiceClient) receiveMessage(ctx context.Context) {
	stream, err := c.client.ReceiveMessage(ctx, &messenger.User{Name: c.name})
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		<-stream.Context().Done()
		fmt.Println(prefixServiceMessage + "Connection to server has lost")
	}()
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			return
		}
		fmt.Printf("%v: %v\n", msg.Sender, msg.Message)
	}
}

func (c *MessengerServiceClient) sendMessage(ctx context.Context, message string) error {
	messageSplit := strings.SplitN(message, " ", 2)
	if len(messageSplit) < 2 {
		fmt.Println(prefixServiceMessage + "Incorrect message. Message should look like \"{username} {message}\"")
		return nil
	}
	response, err := c.client.SendMessage(ctx, &messenger.Message{
		Sender:   c.name,
		Receiver: messageSplit[0],
		Message:  messageSplit[1],
	})
	if err != nil {
		return err
	}
	if !response.Sent {
		fmt.Printf(prefixServiceMessage+"User %s is offline!\n", messageSplit[0])
	}
	return nil
}

func (c *MessengerServiceClient) connect(ctx context.Context) error {
	stream, err := c.client.Ping(ctx)
	if err != nil {
		return err
	}
	err = stream.Send(&messenger.User{Name: c.name})
	if err != nil {
		return err
	}
	_, err = stream.Recv()
	if err != nil {
		return err
	}
	go c.receiveMessage(ctx)
	return nil
}

func (c *MessengerServiceClient) auth(ctx context.Context, scanner *bufio.Scanner) error {
BEGIN:
	fmt.Print("Are you already have a account? ")
	if scanner.Scan() {
		text := strings.ToLower(scanner.Text())
		switch text {
		case "n", "no", "н", "нет":
			fmt.Println(prefixServiceMessage + "SIGN UP")
			if err := c.signUp(ctx, scanner); err != nil {
				return err
			}
			fallthrough
		case "y", "yes", "д", "да":
			fmt.Println(prefixServiceMessage + "LOG IN")
			if err := c.logIn(ctx, scanner); err != nil {
				return err
			}
		default:
			goto BEGIN
		}
	}
	return nil
}

func (c *MessengerServiceClient) signUp(ctx context.Context, scanner *bufio.Scanner) error {
BEGIN:
	var authData model.AuthData
	fmt.Print("Login: ")
	if scanner.Scan() {
		authData.Login = scanner.Text()
	}
	fmt.Print("Name: ")
	if scanner.Scan() {
		authData.Name = scanner.Text()
	}
	for {
		fmt.Print("Password: ")
		if scanner.Scan() {
			authData.Password = scanner.Text()
		}
		fmt.Print("Password: ")
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
		fmt.Println(err)
		goto BEGIN
	}
	return nil
}

func (c *MessengerServiceClient) logIn(ctx context.Context, scanner *bufio.Scanner) error {
BEGIN:
	var login string
	var password string
	fmt.Print("Login: ")
	if scanner.Scan() {
		login = scanner.Text()
	}
	fmt.Print("Password: ")
	if scanner.Scan() {
		password = scanner.Text()
	}
	user, err := c.client.LogIn(ctx, &messenger.LogInRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		fmt.Println(err)
		goto BEGIN
	}
	c.name = user.Name
	return nil
}
