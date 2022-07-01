run: 
	go run main.go

build: 
	go build main.go

test: 
	go test -v ./...

lint: 
	@echo "Applying linter"
	golangci-lint version
	golangci-lint run -c .golangci.yaml ./...
