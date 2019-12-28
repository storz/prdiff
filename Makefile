.PHONY: run
run:
	go run $$(ls -1 cmd/*.go | grep -v _test.go)

.PHONY: build
build:
	go build -o prdiff cmd/*.go
