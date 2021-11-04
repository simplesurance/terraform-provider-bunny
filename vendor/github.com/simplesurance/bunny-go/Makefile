default: build

.PHONY: build
build:
	$(info * compiling)
	go build ./...

.PHONY: check
check:
	$(info * running golangci-lint code checks)
	golangci-lint run

.PHONY: test
test:
	$(info * running tests)
	go test -race -v ./...
