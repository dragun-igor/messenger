@echo off
start "Server" go run cmd/server/main.go
start "Client Igor" go run cmd/client/main.go -name Igor
start "Client Ira" go run cmd/client/main.go -name Ira