@echo off
protoc messenger.proto --go_out=. --go-grpc_out=.