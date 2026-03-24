.PHONY: benchmark
benchmark:
	go test -bench . -benchmem ./...

.PHONY: fuzz
fuzz:
	@grep -rh "^func Fuzz" *_test.go | sed 's/func \(Fuzz[a-zA-Z0-9_]*\).*/\1/' | while read name; do \
		go test -fuzz "^$$name$$" -fuzztime 10s .; \
	done

.PHONY: cover
cover:
	go tool cover -html=cover.out

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -coverprofile=cover.out -shuffle on ./...
