.PHONY: run build test tidy clean

run:
	go run ./cmd/gateway

build:
	go build -o bin/gateway ./cmd/gateway

# test:
# 	go test ./...

tidy:
	go mod tidy

clean:
	rm -rf bin/