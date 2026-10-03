.PHONY: dev dev-backend dev-frontend install build clean help

# Default target
.DEFAULT_GOAL := dev

# Run both backend (Air) and frontend (Vite) concurrently with hot-reload
dev:
	npx --yes concurrently -k -p "[{name}]" -n "BACKEND,FRONTEND" -c "cyan.bold,emerald.bold" "cd backend && air" "cd frontend && npm run dev"

# Run only the backend with Air hot-reloading
dev-backend:
	cd backend && air

# Run only the frontend Vite development server
dev-frontend:
	cd frontend && npm run dev

# Install dependencies for both backend and frontend
install:
	cd backend && go mod download
	cd frontend && npm install

# Build production artifacts for both backend and frontend
build:
	cd backend && go build -o bin/server.exe .
	cd frontend && npm run build

# Clean temporary build artifacts
clean:
	@rm -rf backend/tmp backend/bin frontend/dist 2>/dev/null || true
