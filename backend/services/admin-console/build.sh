#!/bin/bash

# Admin Console Microservice Build Script

set -e

echo "🚀 Building Admin Console Microservice..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    print_error "go.mod not found. Please run this script from the admin-console directory."
    exit 1
fi

# Build Angular application
print_status "Building Angular application..."
cd ../../../adminConsole

if [ ! -f "package.json" ]; then
    print_error "Angular package.json not found. Please ensure the adminConsole directory exists."
    exit 1
fi

# Install dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    print_status "Installing Angular dependencies..."
    npm install
fi

# Build Angular application
print_status "Building Angular application for production..."
npm run build

if [ $? -eq 0 ]; then
    print_status "Angular build completed successfully!"
else
    print_error "Angular build failed!"
    exit 1
fi

# Copy Angular build to admin-console service
print_status "Copying Angular build to admin-console service..."
cd ../backend/services/admin-console

# Create dist directory if it doesn't exist
mkdir -p dist

# Copy Angular build files
cp -r ../../../adminConsole/dist/admin-console/* dist/

# Build Go application
print_status "Building Go application..."
go mod tidy
go build -o admin-console ./cmd/main.go

if [ $? -eq 0 ]; then
    print_status "Go build completed successfully!"
else
    print_error "Go build failed!"
    exit 1
fi

print_status "🎉 Admin Console Microservice build completed successfully!"
print_status "You can now run the service with: ./admin-console"
print_status "Or use Docker: docker build -t scopeapi-admin-console ." 