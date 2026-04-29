.PHONY: build test install clean

BINARY_NAME=vsc

build:
	go build -o $(BINARY_NAME) .

test:
	go test -v ./...

install:
	go install

clean:
	rm -f $(BINARY_NAME)