#!/bin/bash

# Set variables
DEPLOY_DIR="deploy"
SOURCE_DIR="."

# Create fresh deployment directory
rm -rf $DEPLOY_DIR
mkdir -p $DEPLOY_DIR

# Function to create directory if it doesn't exist
create_dir() {
    if [ ! -d "$DEPLOY_DIR/$1" ]; then
        mkdir -p "$DEPLOY_DIR/$1"
    fi
}

echo "🚀 Starting build process..."

make build

# Create necessary directories
create_dir "bin"
create_dir "platform"
create_dir "logs"
create_dir "web"

# Copy essential files and directories
echo "📁 Copying essential files..."

# Copy Docker files
cp Dockerfile "$DEPLOY_DIR/"
cp docker-compose.yml "$DEPLOY_DIR/"

# Copy binary and platform files
cp -r bin/* "$DEPLOY_DIR/bin/" 2>/dev/null || echo "⚠️  No binaries found in bin directory"
cp -r platform/* "$DEPLOY_DIR/platform/"
cp -r web/* "$DEPLOY_DIR/web/"

# Copy config files
cp .gostarter.yaml.docker "$DEPLOY_DIR/.gostarter.yaml"

# Remove any development or temporary files
find "$DEPLOY_DIR" -name "*_test.go" -type f -delete
find "$DEPLOY_DIR" -name ".DS_Store" -type f -delete
find "$DEPLOY_DIR" -name "*.log" -type f -delete
find "$DEPLOY_DIR" -name "*.tmp" -type f -delete

echo "✨ Build complete! Deployment files are ready in the '$DEPLOY_DIR' directory"
echo "📦 You can now use these files for deployment"
