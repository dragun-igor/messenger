@echo off
start "Server" go run cmd/server/main.go
start "Client" go run cmd/client/main.go -pport=9094
start "Client" go run cmd/client/main.go -pport=9095