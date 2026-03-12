VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o clwatch ./cmd/clwatch

install:
	go build $(LDFLAGS) -o $(GOPATH)/bin/clwatch ./cmd/clwatch

clean:
	rm -f clwatch

.PHONY: build install clean
