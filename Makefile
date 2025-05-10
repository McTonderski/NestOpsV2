# Paths
FRONTEND_DIR=frontend
BACKEND_DIR=backend
STATIC_DIR=$(BACKEND_DIR)/static

# Default target
all: build-frontend run-backend

# Build the frontend (assumes Vite)
build-frontend:
	cd $(FRONTEND_DIR) && npm install && npm run build
	mkdir -p $(STATIC_DIR)
	cp -r $(FRONTEND_DIR)/dist/* $(STATIC_DIR)/

# Run the backend
run-backend:
	cd $(BACKEND_DIR) && go run main.go

# Clean static files
clean:
	rm -rf $(STATIC_DIR)/*

# Full rebuild
rebuild: clean all