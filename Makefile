APP_NAME = sleego
BUILD_DIR = ../../bin

GUI_DIR = ./cmd/gui

ICON = ./assets/sleego.ico

ifeq ($(OS),Windows_NT)
	SHELL := cmd
	.SHELLFLAGS := /C
endif
.PHONY: fyne_deps cli linux_gui windows_gui clean

debian_deps:
	apt-get install libgl1-mesa-dev xorg-dev libxkbcommon-dev


fyne_deps_linux:
	@if command -v fyne >/dev/null 2>&1; then \
		echo "Fyne dependencies already installed."; \
	else \
		echo "Installing Fyne dependencies..."; \
		go install fyne.io/tools/cmd/fyne@1.7.0; \
		export PATH=$$PATH:$(go env GOPATH)/bin; \
	fi

fyne_deps_windows:
	@where fyne >nul 2>nul && (echo Fyne dependencies already installed.) || ( \
		echo Installing Fyne dependencies for Windows... & \
		go install fyne.io/tools/cmd/fyne@1.7.0; \
		set "PATH=%PATH%;$(shell go env GOPATH)/bin" \
	)

cli: 
	@echo "Compiling CLI version..."
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/cli



linux_gui_bin: fyne_deps_linux
	@if command -v $(APP_NAME) >/dev/null 2>&1; then \
		printf "%s already installed. Aborting.\n" "$(APP_NAME)" >&2; \
		exit 1; \
	fi
	@echo "Compiling GUI version for Linux..."
	fyne package -os linux -icon $(ICON) -name $(APP_NAME) -executable $(APP_NAME) -source-dir $(GUI_DIR)
	mkdir -p $(APP_NAME)_gui && tar -xf $(APP_NAME).tar.xz -C $(APP_NAME)_gui && rm -rf $(APP_NAME).tar.xz
	make -C $(APP_NAME)_gui user-install
	mkdir -p $(HOME)/.config/$(APP_NAME)
	cp config.json $(HOME)/.config/$(APP_NAME)/config.json
	
linux_gui_remove:
	@echo "Removing GUI version for Linux..."
	make -C $(APP_NAME)_gui user-uninstall
	rm -rf $(HOME)/.config/$(APP_NAME)
	rm -rf $(APP_NAME)_gui
	@echo "Removing GUI version for Linux... done."

linux_gui_rebin:
	@echo "Reinstalling GUI version for Linux..." && \
	$(MAKE) linux_gui_remove && \
	$(MAKE) linux_gui_bin

sleego_debug:
	sleego -loglevel=debug


windows_gui_bin: fyne_deps_windows 
	@echo "Compiling GUI version for Windows..."
	fyne package -os windows -icon $(ICON) -name $(APP_NAME)_gui -executable $(BUILD_DIR)/$(APP_NAME)_gui.exe -sourceDir $(GUI_DIR)
	@if not exist "%APPDATA%\$(APP_NAME)" mkdir "%APPDATA%\$(APP_NAME)"
	@copy /Y config.json "%APPDATA%\$(APP_NAME)\config.json"

clean:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)

test:
	@echo "Running tests..."
	go test -v ./...

govuln:
	@echo "Running govulncheck..."
	govulncheck -show verbose ./... 
	