package test

import (
	"context"
	"testing"

	"github.com/dragun-igor/messenger/proto/messenger"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestClientFlow(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	ctx := context.Background()
	messengerService := messenger.NewMessengerServiceClient(conn)

	// sign up
	signUpRequest := &messenger.SignUpRequest{
		Login:    "Login",
		Name:     "Name",
		Password: "Password",
	}
	_, err = messengerService.SignUp(ctx, signUpRequest)
	require.NoError(t, err)
	_, err = messengerService.SignUp(ctx, signUpRequest)
	require.EqualError(t, err, "rpc error: code = AlreadyExists desc = login name is busy")
	signUpRequest.Login = "EditedLogin"
	_, err = messengerService.SignUp(ctx, signUpRequest)
	require.EqualError(t, err, "rpc error: code = AlreadyExists desc = user name is busy")

	// log in
	logInRequest := &messenger.LogInRequest{
		Login:    "Login",
		Password: "WrongPassword",
	}
	_, err = messengerService.LogIn(ctx, logInRequest)
	require.EqualError(t, err, "rpc error: code = PermissionDenied desc = incorrect password")
	logInRequest.Password = "Password"
	user, err := messengerService.LogIn(ctx, logInRequest)
	require.NoError(t, err)
	require.Equal(t, signUpRequest.Name, user.Name)

	// send message
	message := &messenger.Message{
		Sender:   "Name",
		Receiver: "User",
		Message:  "Hello!",
	}
	resp, err := messengerService.SendMessage(ctx, message)
	require.NoError(t, err)
	require.Equal(t, false, resp.Sent)
}
