APP_NAME = sleego
BUILD_DIR = ../../bin

GUI_DIR = ./cmd/gui

ICON = ./assets/sleego_icon.png

.PHONY: fyne_deps cli linux_gui windows_gui clean

debian_deps:
	sudo apt-get install libgl1-mesa-dev xorg-dev libxkbcommon-dev

fyne_deps:
	@if command -v fyne >/dev/null 2>&1; then \
		echo "Fyne dependencies already installed."; \
	else \
		echo "Instalando Fyne dependencies..."; \
		go get fyne.io/fyne/v2; \
		go install fyne.io/fyne/v2/cmd/fyne@latest; \
		export PATH=$$PATH:$(go env GOPATH)/bin; \
	fi

cli: 
	@echo "Compiling CLI version..."
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/cli



linux_gui: fyne_deps
	@if command -v $(APP_NAME) >/dev/null 2>&1; then \
		printf "%s already installed. Aborting.\n" "$(APP_NAME)" >&2; \
		exit 1; \
	fi
	
	@echo "Compiling GUI version for Linux..."
	fyne package -os linux -icon $(ICON) -name $(APP_NAME)_gui -executable $(APP_NAME) -sourceDir $(GUI_DIR)
	mkdir -p $(APP_NAME)_gui && tar -xf $(APP_NAME)_gui.tar.xz -C $(APP_NAME)_gui && rm -rf $(APP_NAME)_gui.tar.xz
	sudo make -C $(APP_NAME)_gui install
	mkdir -p $(HOME)/.config/$(APP_NAME)
	cp config.json $(HOME)/.config/$(APP_NAME)/config.json
	

windows_gui: fyne_deps 
	@echo "Compiling GUI version for Windows..."
	fyne package -os windows -icon $(ICON) -name $(APP_NAME)_gui -executable $(BUILD_DIR)/$(APP_NAME)_gui.exe -sourceDir $(GUI_DIR)

clean:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)
