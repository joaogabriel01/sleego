APP_NAME = sleego
BUILD_DIR = ../../bin


debian_deps:
	apt-get install libgl1-mesa-dev xorg-dev libxkbcommon-dev


cli: 
	@echo "Compiling CLI version..."
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/cli


clean:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)

test:
	@echo "Running tests..."
	go test -v ./...

govuln:
	@echo "Running govulncheck..."
	govulncheck -show verbose ./... 
	