#!/bin/bash

# Soma Mayel Campaign Website Startup Script

echo "Starting Soma Mayel Campaign Website..."
echo "========================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Load environment variables
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
    echo "✓ Environment variables loaded"
else
    echo "⚠ Warning: .env file not found. Using defaults."
    cp .env.example .env 2>/dev/null || true
fi

# Create necessary directories
mkdir -p content/posts content/pages static/images static/videos
echo "✓ Directories initialized"

# Install dependencies
echo "Installing dependencies..."
go mod download
echo "✓ Dependencies installed"

# Build the application
echo "Building application..."
go build -o bin/soma-campaign main.go
if [ $? -eq 0 ]; then
    echo "✓ Build successful"
else
    echo "✗ Build failed"
    exit 1
fi

# Start the application
echo ""
echo "========================================="
echo "Starting server on http://localhost:${PORT:-3000}"
echo "Press Ctrl+C to stop"
echo "========================================="
echo ""

./bin/soma-campaign