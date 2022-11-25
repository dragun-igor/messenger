# Messenger Service

Console GRPC messaging service.

## Run server

Install protobuf: `sudo apt install -y protobuf-compiler`

Install mockgen `go install github.com/golang/mock/mockgen@latest`

Generate protobuf: `make generate-proto`

Generate mocks: `make generate-mocks`

Run app with environment: `make deploy`

## Run client

Run client `go run cmd/client/main.go`

After authorization user be able send message to other authorizated users. Message should look like `{username} {message_text}`.
