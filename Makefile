all: build

build:
	go build -o termi ./cmd/termi
	go build *.go

clean:
	rm -f termi

run:
	go run ./cmd/termi

fmt:
	go fmt ./...

test:
	go test ./...
