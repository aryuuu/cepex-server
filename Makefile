run: 
	export $$(xargs < .env) && \
	go run main.go

build: 
	go build main.go

test: 
	export $$(xargs < .env) && \
	go test -v ./...

dev-air:
	export $$(xargs < .env) && \
	air

lint: 
	@echo "Applying linter"
	golangci-lint version
	golangci-lint run -c .golangci.yaml ./...
