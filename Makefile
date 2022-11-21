.PHONY: generate-proto
generate-proto:
	@cd proto && make generate

.PHONY: test
test:
	go test -cover --tags=ci ./...

.PHONY: mocks
mocks:
	mockgen -package=mocks -source internal/server/resources/interface.go -destination internal/server/service/mocks/mock.go