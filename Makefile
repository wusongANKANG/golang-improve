MODULE ?= ./examples/02_basics
TEST ?= .
BENCH ?= .

.PHONY: list test module run bench race cover

list:
	go list ./examples/...

test:
	go test ./...

module:
	go test -v $(MODULE)

run:
	go test -run $(TEST) -v $(MODULE)

bench:
	go test -bench $(BENCH) -benchmem $(MODULE)

race:
	go test -race $(MODULE)

cover:
	go test -cover ./...
