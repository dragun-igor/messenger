# Messenger Service

Console GRPC messaging service.

## Run

Install protobuf: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

Install mockgen `go install github.com/golang/mock/mockgen@latest`

Generate protobuf: `make generate-proto`

Generate mocks: `make generate-mocks`

Run app with environment: `make deploy`
