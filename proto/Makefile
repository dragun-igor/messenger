generate:
	@cd messenger && rm -rf *
	@protoc *.proto  -I . -I .. --go_out=. --go-grpc_out=.