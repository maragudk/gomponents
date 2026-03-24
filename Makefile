.PHONY: benchmark
benchmark:
	go test -bench . -benchmem ./...

.PHONY: fuzz
fuzz:
	go test -fuzz FuzzEl -fuzztime 10s .
	go test -fuzz FuzzAttr -fuzztime 10s .
	go test -fuzz FuzzText -fuzztime 10s .
	go test -fuzz FuzzRaw -fuzztime 10s .

.PHONY: cover
cover:
	go tool cover -html=cover.out

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -coverprofile=cover.out -shuffle on ./...
