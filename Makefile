GO=go
BINARY=paxossim
# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

info:
	@echo "---------------------------------------"
	@echo
	@echo " build   generate a build              "
	@echo " test    run unit-tests                "
	@echo " fmt     format using go fmt           "
	@echo "---------------------------------------" 

build: clean fmt
	$(GO) build -o bin/$(BINARY) -v cmd/cmd.go

test: clean fmt
	$(GO) test -v 

fmt:
	$(GO) fmt ./...

clean:
	rm -f bin/$(BINARY)
