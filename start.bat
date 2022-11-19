@echo off
start "Server" go run cmd/server/main.go
start "Client" go run cmd/client/main.go
start "Client" go run cmd/client/main.go