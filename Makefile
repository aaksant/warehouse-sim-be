.PHONY: run test vet tidy

tidy:
	go mod tidy

run:
	go run ./cmd/api

test:
	go test ./... -v

vet:
	go vet ./...
