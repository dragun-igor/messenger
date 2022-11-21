generate:
	@cd messenger && rm -rf *
	@protoc *.proto  -I . -I .. --go_out=plugins=grpc:messenger
	@echo "Client generated"