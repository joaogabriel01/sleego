APP_NAME = sleego
BUILD_DIR = ../../bin

GUI_DIR = ./cmd/gui

ICON = ./assets/sleego_icon.png

.PHONY: fyne_deps cli linux_gui windows_gui clean

fyne_deps:
	@echo "Installing Fyne dependencies..."
	go get fyne.io/fyne/v2
	go install fyne.io/fyne/v2/cmd/fyne@latest
	export PATH=$PATH:$(go env GOPATH)/bin

cli: 
	@echo "Compiling CLI version..."
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/cli

linux_gui: fyne_deps
	sudo apt-get install libgl1-mesa-dev xorg-dev libxkbcommon-dev
	@echo "Compiling GUI version for Linux..."
	fyne package -os linux -icon $(ICON) -name $(APP_NAME)_gui -sourceDir ./cmd/gui
	mkdir -p sleego_gui && tar -xf sleego_gui.tar.xz -C sleego_gui && rm -rf sleego_gui.tar.xz
	sudo make -C sleego_gui install
	mkdir -p $(HOME)/.config/sleego
	cp config.json $(HOME)/.config/sleego/config.json
	

windows_gui: fyne_deps 
	@echo "Compiling GUI version for Windows..."
	fyne package -os windows -icon $(ICON) -name $(APP_NAME)_gui -executable $(BUILD_DIR)/$(APP_NAME)_gui.exe -sourceDir $(GUI_DIR)

clean:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)
