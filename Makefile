.PHONY: benchmark cover lint test

benchmark:
	go test -bench=.

cover:
	go tool cover -html=cover.out

lint:
	golangci-lint run

test:
	go test -coverprofile=cover.out -mod=readonly ./...

