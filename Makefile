.PHONY: cover lint test

cover:
	go tool cover -html=cover.out

lint:
	golangci-lint run

test:
	go test -coverprofile=cover.out -mod=readonly ./...

