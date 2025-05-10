# Makefile for building cross-platform installer binaries

APP_NAME=installer
SRC=installer.go
BUILD_DIR=build

all: linux macos windows

linux:
	@echo "🔨 Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux $(SRC)

macos:
	@echo "🔨 Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin $(SRC)

windows:
	@echo "🔨 Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME).exe $(SRC)

clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

prepare:
	@mkdir -p $(BUILD_DIR)

.PHONY: all linux macos windows clean prepare