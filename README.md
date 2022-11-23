# Messenger Service

Console GRPC messaging service.

## Run

Install protobuf: `sudo snap install protobuf --classic`

Install mockgen `go install github.com/golang/mock/mockgen@latest`

Generate protobuf: `make generate-proto`

Generate mocks: `make generate-mocks`

Run app with environment:<br/>
`cd deployments`<br/>
`sudo docker-compose up --build`
