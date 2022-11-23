.PHONY: generate-proto
generate-proto:
	@cd proto && make generate

.PHONY: test
test:
	go test -cover --tags=ci ./...

.PHONY: generate-mocks
generate-mocks:
	go generate internal/server/service/interface.go

.PHONY: deploy
deploy:
	@cd deployments && sudo docker-compose up --build