@echo off
cd c:/go/messenger
protoc messengerpb/messenger.proto --go_out=. --go-grpc_out=.
go build server/server.go
go build client/client.go
start "Server" "server.exe"
start "Client" "client.exe"
start "Client" "client.exe"