# Messenger Service

Console GRPC messaging service.

## Run

Install protobuf: `sudo apt install -y protobuf-compiler`

Install mockgen `go install github.com/golang/mock/mockgen@latest`

Generate protobuf: `make generate-proto`

Generate mocks: `make generate-mocks`

Run app with environment: `make deploy`
