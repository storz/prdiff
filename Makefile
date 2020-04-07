.PHONY: run
run:
	go run $$(ls -1 cmd/prdiff/*.go | grep -v _test.go)

.PHONY: build
build:
	go build ./cmd/prdiff

.PHONY: test
test:
	go test -v -race -count=1 ./...
