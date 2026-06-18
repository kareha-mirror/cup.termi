all: build

build:
	go build -o termi ./cmd/termi
	go build

clean:
	rm -f termi

run:
	go run ./cmd/termi

fmt:
	go fmt ./...

test:
	go test ./...

windows:
	GOOS=windows go build -o termi-windows ./cmd/termi
	GOOS=windows go build
