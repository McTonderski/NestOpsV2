# Paths
FRONTEND_DIR=frontend
BACKEND_DIR=backend
STATIC_DIR=$(BACKEND_DIR)/static
DIST_DIR=dist
ARCHIVE_NAME=homelab.tar.gz

# Default target
all: clean build-frontend build-backend package

# Build the frontend
build-frontend:
	cd $(FRONTEND_DIR) && npm install && npm run build
	mkdir -p $(STATIC_DIR)
	cp -r $(FRONTEND_DIR)/dist/* $(STATIC_DIR)/

# Build the backend binary
build-backend:
	GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/homelab $(BACKEND_DIR)/main.go

# Package everything into a tar.gz archive
package:
	mkdir -p $(DIST_DIR)/static
	cp -r $(STATIC_DIR)/* $(DIST_DIR)/static/
	tar -czf $(ARCHIVE_NAME) -C $(DIST_DIR) .

# Clean all build artifacts
clean:
	rm -rf $(DIST_DIR) $(ARCHIVE_NAME)

.PHONY: all build-frontend build-backend package clean