run: 
	go run main.go

build: 
	go build main.go

lint: 
	@echo "Applying linter"
	golangci-lint version
	golangci-lint run -c .golangci.yaml ./...
