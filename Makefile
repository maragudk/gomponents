.PHONY: benchmark
benchmark:
	go test -bench . -benchmem ./...

.PHONY: fuzz
fuzz:
	grep -r "^func Fuzz" --include='*_test.go' -l . | while read file; do \
		dir=$$(dirname "$$file"); \
		grep -h "^func Fuzz" "$$file" | sed 's/func \(Fuzz[a-zA-Z0-9_]*\).*/\1/' | while read name; do \
			go test -fuzz "^$$name$$" -fuzztime 10s "$$dir" || exit 1; \
		done || exit 1; \
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
