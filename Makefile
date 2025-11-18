# Project configuration
BINARY_NAME := tgs

EXTENSION_TAGS := stats,ip,journal,opmanga,pgpress

# Directories and paths
BIN_DIR := ./bin

# Service configuration
LAUNCHD_SERVICE_FILE := ~/Library/LaunchAgents/com.$(BINARY_NAME).plist

# Logs
SERVICE_LOG := $(HOME)/Library/Application Support/$(BINARY_NAME).log

.PHONY: all clean init run build macos-install macos-start-service \
	 macos-stop-service macos-print-service macos-watch-service

all: init build

clean:
	git clean -xfd

init:
	go mod tidy -v

run:
	go run --tags=${EXTENSION_TAGS} -v ./cmd/tgs

build:
	go build --tags=${EXTENSION_TAGS} -v -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/tgs

# TODO: ...

# Create launchd service file for macOS
define LAUNCHD_SERVICE_FILE_CONTENT
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.$(BINARY_NAME)</string>

	<key>ProgramArguments</key>
	<array>
		<string>/usr/local/bin/$(BINARY_NAME)</string>
	</array>

	<key>RunAtLoad</key>
	<true/>

	<key>KeepAlive</key>
	<true/>

	<key>StandardOutPath</key>
	<string>$(SERVICE_LOG)</string>

	<key>StandardErrorPath</key>
	<string>$(SERVICE_LOG)</string>
</dict>
</plist>
endef

export LAUNCHD_SERVICE_FILE_CONTENT

macos-install:
	@echo "Installing $(BINARY_NAME) for macOS..."
	mkdir -p /usr/local/bin
	sudo cp $(BIN_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "$$LAUNCHD_SERVICE_FILE_CONTENT" > $(LAUNCHD_SERVICE_FILE)
	@echo "$(BINARY_NAME) installed successfully"

macos-start-service:
	@echo "Starting $(BINARY_NAME) service..."
	launchctl load -w $(LAUNCHD_SERVICE_FILE)
	launchctl start com.$(BINARY_NAME)

macos-stop-service:
	@echo "Stopping $(BINARY_NAME) service..."
	launchctl stop com.$(BINARY_NAME)
	launchctl unload -w $(LAUNCHD_SERVICE_FILE)

macos-restart-service:
	@echo "Restarting $(BINARY_NAME) service..."
	make macos-stop-service
	make macos-start-service

macos-print-service:
	@echo "$(BINARY_NAME) service information:"
	@launchctl print gui/$$(id -u)/com.$(BINARY_NAME) || echo "Service not loaded or running"

macos-watch-service:
	@echo "$(BINARY_NAME) watch server logs @ \"$(SERVICE_LOG)\":"
	@if [ -f "$(SERVICE_LOG)" ]; then \
		echo "Watching logs... Press Ctrl+C to stop"; \
		tail -f "$(SERVICE_LOG)"; \
	else \
		echo "Log file not found. Make sure the service is running or has been started."; \
		echo "Log file path: \"$(SERVICE_LOG)\""; \
	fi
