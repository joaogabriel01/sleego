APP_NAME = sleego

cli: 
	@echo "Compiling CLI version..."
	go build -o $(APP_NAME) ./cmd/cli

test:
	@echo "Running tests..."
	go test -v ./...
