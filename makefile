.PHONY: docs
docs:
	swag init -g cmd/api/main.go 

.PHONY: run
run:
	go run ./cmd/api

.PHONY: mocks
mocks:
	go generate ./..

.PHONY: test
test:
	go test ./... -v
